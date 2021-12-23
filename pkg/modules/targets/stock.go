// Package targets package targets
package targets

import (
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

// getStockTargets getStockTargets
func getStockTargets() error {
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
		eventbus.Get().Pub(eventbus.TopicTargets(), inDBTargets)
		return nil
	}

	lastTradeDayArr := tradeday.GetLastNTradeDayByDate(1, tradeDay)
	err = sinopacapi.Get().FetchVolumeRankByDate(lastTradeDayArr[0].Format(global.ShortTimeLayout), 200)
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

	conf, err := config.Get()
	if err != nil {
		log.Get().Panic(err)
	}
	condition := conf.GetTradeTargetCondtion()

	var targetArr []*dbagent.Target
	for _, v := range body.GetData() {
		if stockTargetFilter(v, condition) {
			tmp := &dbagent.Target{
				Stock:    cache.GetCache().GetStock(v.GetCode()),
				TradeDay: cache.GetCache().GetTradeDay(),
				Volume:   v.GetTotalVolume(),
				Rank:     len(targetArr) + 1,
			}
			targetArr = append(targetArr, tmp)
			log.Get().WithFields(map[string]interface{}{
				"Stock":       tmp.Stock.Name,
				"TotalVolume": tmp.Volume,
				"Rank":        tmp.Rank,
			}).Info("Target")
		}
	}

	err = dbagent.Get().DeleteMultiTargetByDate(cache.GetCache().GetTradeDay())
	if err != nil {
		log.Get().Panic(err)
	}
	err = dbagent.Get().InsertMultiTarget(targetArr)
	if err != nil {
		log.Get().Panic(err)
	}

	// send to bus
	eventbus.Get().Pub(eventbus.TopicTargets(), targetArr)
}

func stockTargetFilter(v *pb.VolumeRankMessage, cond *config.TargetCond) bool {
	if v.GetTotalVolume() < cond.LimitVolume {
		return false
	}
	if v.GetClose() < cond.LimitPriceLow || v.GetClose() > cond.LimitPriceHigh {
		return false
	}
	return true
}
