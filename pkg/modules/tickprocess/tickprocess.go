// Package tickprocess package tickprocess
package tickprocess

import (
	"time"
	"trade_agent/pkg/cache"
	"trade_agent/pkg/config"
	"trade_agent/pkg/dbagent"
	"trade_agent/pkg/eventbus"
	"trade_agent/pkg/log"
	"trade_agent/pkg/sinopacapi"
)

var (
	simTradeRealTimeTickChannel   = make(chan int)
	simTradeRealTimeBidAskChannel = make(chan int)
)

// InitTickProcess InitTickProcess
func InitTickProcess() {
	log.Get().Info("Initial TickProcess")

	go simTradeRealTimeTickCollector()
	go simTradeRealTimeBidAskCollector()

	// sub targets to sub mq history tick, kbar, realtime tick, bidask
	err := eventbus.Get().Sub(eventbus.TopicTargets(), targetsBusCallback)
	if err != nil {
		log.Get().Panic(err)
	}

	err = subHistroyTick()
	if err != nil {
		log.Get().Panic(err)
	}
	err = subHistroyKbar()
	if err != nil {
		log.Get().Panic(err)
	}
}

func targetsBusCallback(targetArr []*dbagent.Target) error {
	for _, v := range targetArr {
		cache.GetCache().SetRealTimeTickChannel(v.Stock.Number, make(chan *dbagent.RealTimeTick))
		go realTimeTickProcessor(v.Stock.Number)

		cache.GetCache().SetRealTimeBidAskChannel(v.Stock.Number, make(chan *dbagent.RealTimeBidAsk))
		go realTimeBidAskProcessor(v.Stock.Number)
	}

	err := subRealTimeTick()
	if err != nil {
		log.Get().Panic(err)
	}
	err = subRealTimeBidAsk()
	if err != nil {
		log.Get().Panic(err)
	}

	eventbus.Get().Pub(eventbus.TopicSubscribeTargets(), targetArr)
	return nil
}

func realTimeTickProcessor(stockNum string) {
	analyzeConf := config.GetAnalyzeConfig()
	ch := cache.GetCache().GetRealTimeTickChannel(stockNum)

	var tickArr dbagent.RealTimeTickArr
	var lastPeriodEndTime time.Time
	for {
		tick := <-ch
		tickArr = append(tickArr, tick)

		// save realtime tick close to cache
		cache.GetCache().SetRealTimeTickClose(stockNum, tick.Close)

		if lastPeriodEndTime.IsZero() {
			lastPeriodEndTime = tick.TickTime
			continue
		}

		lastPeriodArr := tickArr.GetLastNSecondArr(analyzeConf.TickAnalyzeMinPeriod)
		if float64(len(lastPeriodArr)) < analyzeConf.MaxLoss {
			continue
		}

		if tick.TickTime.Before(lastPeriodEndTime.Add(time.Duration(analyzeConf.TickAnalyzeMinPeriod) * time.Second)) {
			continue
		} else {
			lastPeriodEndTime = lastPeriodArr.GetLastTick().TickTime
		}

		realTimeBalancePct, emerAction := getRealTimeBalancePct(stockNum, tick.Close)
		action := realTimeTickArrActionGenerator(tickArr, lastPeriodArr, analyzeConf)
		if action == 0 && realTimeBalancePct < analyzeConf.MaxLoss {
			continue
		} else {
			action = emerAction
		}

		order := &sinopacapi.Order{
			StockNum:  stockNum,
			Price:     tick.Close,
			Action:    action,
			TradeTime: tick.TickTime,
		}

		// send order event
		eventbus.Get().Pub(eventbus.TopicStockOrder(), order)
	}
}

func realTimeBidAskProcessor(stockNum string) {
	ch := cache.GetCache().GetRealTimeBidAskChannel(stockNum)
	var bidAskArr []*dbagent.RealTimeBidAsk
	for {
		bidAsk := <-ch
		bidAskArr = append(bidAskArr, bidAsk)
		status := realTimeBidAskArrStatusGenerator(bidAsk, bidAskArr)
		cache.GetCache().SetRealTimeBidAskStatus(stockNum, status)
	}
}

func simTradeRealTimeTickCollector() {
	printMinute := time.Now().Minute()
	var count int
	for {
		simTrade := <-simTradeRealTimeTickChannel
		count += simTrade
		if time.Now().Minute() != printMinute {
			printMinute = time.Now().Minute()
			log.Get().WithFields(map[string]interface{}{
				"Count": count,
			}).Info("SimTradeTick")
		}
	}
}

func simTradeRealTimeBidAskCollector() {
	printMinute := time.Now().Minute()
	var count int
	for {
		simTrade := <-simTradeRealTimeBidAskChannel
		count += simTrade
		if time.Now().Minute() != printMinute {
			printMinute = time.Now().Minute()
			log.Get().WithFields(map[string]interface{}{
				"Count": count,
			}).Info("SimTradeBidAsk")
		}
	}
}

func getRealTimeBalancePct(stockNum string, close float64) (float64, sinopacapi.OrderAction) {
	historyOrderBuy := cache.GetCache().GetOrderBuy(stockNum)
	historyOrderSell := cache.GetCache().GetOrderSell(stockNum)
	restOrderCount := len(historyOrderBuy) - len(historyOrderSell)
	if restOrderCount != 0 {
		for i := len(historyOrderSell); i <= len(historyOrderSell)-1+restOrderCount; i++ {
			lastOrder := historyOrderBuy[i]
			buyCost := sinopacapi.GetStockBuyCost(lastOrder.Price, lastOrder.Quantity)
			sellCost := sinopacapi.GetStockSellCost(close, lastOrder.Quantity)
			if buyCost > sellCost {
				return 100 * float64(buyCost-sellCost) / float64(buyCost), sinopacapi.ActionSell
			}
		}
	}

	historyOrderSellFirst := cache.GetCache().GetOrderSellFirst(stockNum)
	historyOrderBuyLater := cache.GetCache().GetOrderBuyLater(stockNum)
	restOrderCount = len(historyOrderSellFirst) - len(historyOrderBuyLater)
	if restOrderCount != 0 {
		for i := len(historyOrderBuyLater); i <= len(historyOrderBuyLater)-1+restOrderCount; i++ {
			lastOrder := historyOrderSellFirst[i]
			sellFirstCost := sinopacapi.GetStockSellCost(lastOrder.Price, lastOrder.Quantity)
			buyLaterCost := sinopacapi.GetStockBuyCost(close, lastOrder.Quantity)
			if sellFirstCost < buyLaterCost {
				return 100 * float64(buyLaterCost-sellFirstCost) / float64(sellFirstCost), sinopacapi.ActionBuyLater
			}
		}
	}
	return 0, 0
}
