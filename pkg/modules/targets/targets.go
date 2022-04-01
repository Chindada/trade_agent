// Package targets package targets
package targets

import (
	"time"

	"trade_agent/global"
	"trade_agent/pkg/cache"
	"trade_agent/pkg/config"
	"trade_agent/pkg/dbagent"
	"trade_agent/pkg/eventbus"
	"trade_agent/pkg/log"
	"trade_agent/pkg/modules/tradeday"
	"trade_agent/pkg/mqhandler"
	"trade_agent/pkg/pb"
	"trade_agent/pkg/sinopacapi"
)

var lastNTradeDayTarget int64 = 1

// InitTargets InitTargets
func InitTargets() {
	log.Get().Info("Initial Targets")

	err := getStockTargets(lastNTradeDayTarget)
	if err != nil {
		log.Get().Panic(err)
	}

	go func() {
		for range time.NewTicker(60 * time.Second).C {
			if cache.GetCache().GetIsAllowTrade() {
				err = getRealTimeTargets()
				if err != nil {
					log.Get().Panic(err)
				}
			}
		}
	}()

	go func() {
		for range time.NewTicker(30 * time.Second).C {
			err = getTSERealTime()
			if err != nil {
				log.Get().Panic(err)
			}
		}
	}()
}

// getStockTargets getStockTargets
func getStockTargets(n int64) error {
	handler := mqhandler.Get()
	err := handler.Sub(mqhandler.MQSubBody{
		MQTopic:  mqhandler.TopicVolumeRank(),
		Once:     true,
		Callback: volumeRankCallback,
	})
	if err != nil {
		return err
	}

	tradeDay := cache.GetCache().GetTradeDay()
	inDBTargets, err := dbagent.Get().GetTargetsByDate(tradeDay)
	if err != nil {
		return err
	}
	if len(inDBTargets) != 0 {
		for _, v := range inDBTargets {
			log.Get().WithFields(map[string]interface{}{
				"Stock":       v.Stock.Name,
				"TotalVolume": v.Volume,
				"Rank":        v.Rank,
			}).Info("DB Target")
		}
		// send to bus
		eventbus.Get().PublishTargets(inDBTargets)
		// append to cache
		cache.GetCache().AppendTargets(inDBTargets)
		err = handler.UnSub(string(mqhandler.TopicVolumeRank()))
		if err != nil {
			return err
		}
		return nil
	}

	lastTradeDay := tradeday.GetLastNTradeDayByDate(n, tradeDay)[n-1]
	err = sinopacapi.Get().FetchVolumeRankByDate(lastTradeDay.Format(global.ShortTimeLayout), 200)
	if err != nil {
		return err
	}
	return nil
}

// volumeRankCallback volumeRankCallback
func volumeRankCallback(m mqhandler.MQMessage) {
	body := pb.VolumeRankResponse{}
	err := body.UnmarshalProto(m.Payload())
	if err != nil {
		log.Get().Panic(err)
	}
	if len(body.GetData()) == 0 {
		lastNTradeDayTarget++
		err = getStockTargets(lastNTradeDayTarget)
		if err != nil {
			log.Get().Panic(err)
		}
		return
	}

	condition := config.GetTargetCondConfig()
	tradeDay := cache.GetCache().GetTradeDay()
	var targetArr []*dbagent.Target
	for _, v := range body.GetData() {
		stock := cache.GetCache().GetStock(v.GetCode())
		if stock == nil {
			log.Get().WithFields(map[string]interface{}{
				"Stock": v.GetCode(),
			}).Error("VolumeRank Stock Cache Error")
			continue
		}
		// save volume to cache
		cache.GetCache().SetStockVolume(v.GetCode(), v.GetTotalVolume())
		tmpTarget := stockWithData{
			stock:       stock,
			close:       v.GetClose(),
			totalVolume: v.GetTotalVolume(),
		}
		if stockTargetFilter(tmpTarget, condition, false) {
			tmp := &dbagent.Target{
				Stock:       tmpTarget.stock,
				TradeDay:    tradeDay,
				Volume:      tmpTarget.totalVolume,
				Rank:        len(targetArr) + 1,
				RealTimeAdd: false,
				Subscribe:   true,
			}
			targetArr = append(targetArr, tmp)
			log.Get().WithFields(map[string]interface{}{
				"Stock":       tmp.Stock.Name,
				"TotalVolume": tmp.Volume,
				"Rank":        tmp.Rank,
			}).Info("Target")
		}
	}

	if len(targetArr) >= 100 {
		targetArr = targetArr[:100]
	}

	err = dbagent.Get().DeleteMultiTargetByDate(tradeDay)
	if err != nil {
		log.Get().Panic(err)
	}
	err = dbagent.Get().InsertMultiTarget(targetArr)
	if err != nil {
		log.Get().Panic(err)
	}

	// send to bus
	eventbus.Get().PublishTargets(targetArr)
	// append to cache
	cache.GetCache().AppendTargets(targetArr)
}
