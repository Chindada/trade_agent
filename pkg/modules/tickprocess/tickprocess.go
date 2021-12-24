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

	log.Get().Info("Initial TickProcess")
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

	return nil
}

func realTimeTickProcessor(stockNum string) {
	ch := cache.GetCache().GetRealTimeTickChannel(stockNum)
	for {
		tick := <-ch
		log.Get().Info(tick.Stock.Name)

		order := sinopacapi.Order{
			StockNum: stockNum,
			Price:    tick.Close,
			Quantity: 1,
			Action:   sinopacapi.ActionBuy,
		}
		eventbus.Get().Pub(eventbus.TopicStockOrder(), order)
	}
}

func realTimeBidAskProcessor(stockNum string) {
	ch := cache.GetCache().GetRealTimeBidAskChannel(stockNum)
	for {
		tick := <-ch

		log.Get().Info(tick.Stock.Name)
	}
}
