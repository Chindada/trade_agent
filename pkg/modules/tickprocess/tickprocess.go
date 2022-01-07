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
		cache.GetCache().Set(cache.KeyRealTimeTickChannel(v.Stock.Number), make(chan *dbagent.RealTimeTick))
		go realTimeTickProcessor(v.Stock.Number)

		cache.GetCache().Set(cache.KeyRealTimeBidAskChannel(v.Stock.Number), make(chan *dbagent.RealTimeBidAsk))
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
	for {
		tick := <-ch
		tickArr = append(tickArr, tick)
		action := realTimeTickArrActionGenerator(tickArr, analyzeConf)
		if action == 0 {
			continue
		}

		order := &sinopacapi.Order{
			StockNum:  stockNum,
			Price:     tick.Close,
			Action:    action,
			TradeTime: tick.TickTime,
		}
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
		cache.GetCache().Set(cache.KeyRealTimeBidAskStatus(stockNum), status)
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
