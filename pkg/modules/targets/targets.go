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

// InitTargets InitTargets
func InitTargets() {
	log.Get().Info("Initial Targets")

	err := getStockTargets()
	if err != nil {
		log.Get().Panic(err)
	}
}

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

	lastTradeDay := tradeday.GetLastNTradeDayByDate(1, tradeDay)[0]
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

	conf, err := config.Get()
	if err != nil {
		log.Get().Panic(err)
	}
	condition := conf.GetTargetCondtion()
	tradeDay := cache.GetCache().GetTradeDay()
	var targetArr []*dbagent.Target
	for _, v := range body.GetData() {
		if stockTargetFilter(v, condition) {
			tmp := &dbagent.Target{
				Stock:    cache.GetCache().GetStock(v.GetCode()),
				TradeDay: tradeDay,
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

	err = dbagent.Get().DeleteMultiTargetByDate(tradeDay)
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

func stockTargetFilter(v *pb.VolumeRankMessage, cond config.TargetCond) bool {
	stock := cache.GetCache().GetStock(v.GetCode())
	if stock == nil {
		log.Get().Errorf("Stock %s does not exist in cache", v.GetCode())
		return false
	}

	blackCategoryMap := make(map[string]bool)
	blackStockMap := make(map[string]bool)
	for _, v := range cond.BlackCategory {
		blackCategoryMap[v] = true
	}
	for _, v := range cond.BlackStock {
		blackStockMap[v] = true
	}

	if blackStockMap[v.GetCode()] {
		return false
	}
	if blackCategoryMap[stock.Category] {
		return false
	}
	if v.GetTotalVolume() < cond.LimitVolume {
		return false
	}
	if v.GetClose() < cond.LimitPriceLow || v.GetClose() > cond.LimitPriceHigh {
		return false
	}
	return true
}
