// Package tickprocess package tickprocess
package tickprocess

import (
	"trade_agent/pkg/cache"
	"trade_agent/pkg/dbagent"
	"trade_agent/pkg/eventbus"
	"trade_agent/pkg/log"
	"trade_agent/pkg/sinopacapi"
)

// InitTickProcess InitTickProcess
func InitTickProcess() {
	log.Get().Info("Initial TickProcess")

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
	ch := cache.GetCache().GetRealTimeTickChannel(stockNum)
	var tickArr []*dbagent.RealTimeTick
	for {
		tick := <-ch
		if tick.Simtrade == 1 {
			continue
		}

		action := realTimeTickArrAnalyzer(tick, tickArr)
		tickArr = append(tickArr, tick)
		if action == 0 {
			continue
		}

		order := &sinopacapi.Order{
			StockNum: stockNum,
			Price:    tick.Close,
			Action:   action,
		}
		eventbus.Get().Pub(eventbus.TopicStockOrder(), order)
	}
}

func realTimeBidAskProcessor(stockNum string) {
	ch := cache.GetCache().GetRealTimeBidAskChannel(stockNum)
	var bidAskArr []*dbagent.RealTimeBidAsk
	for {
		bidAsk := <-ch
		if bidAsk.Simtrade == 1 {
			continue
		}

		bidAskArr = append(bidAskArr, bidAsk)
		status := realTimeBidAskArrAnalyzer(bidAsk, bidAskArr)
		cache.GetCache().Set(cache.KeyRealTimeBidAskStatus(stockNum), status)
	}
}
