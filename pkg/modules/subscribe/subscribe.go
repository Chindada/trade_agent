// Package subscribe package subscribe
package subscribe

import (
	"trade_agent/pkg/dbagent"
	"trade_agent/pkg/eventbus"
	"trade_agent/pkg/log"
	"trade_agent/pkg/sinopacapi"
)

// InitSubscribe InitSubscribe
func InitSubscribe() {
	log.Get().Info("Initial Subscribe")

	// unsubscribe all realtime tick
	err := sinopacapi.Get().UnSubscribeAllByType(sinopacapi.TickTypeStockRealTime)
	if err != nil {
		log.Get().Panic(err)
	}

	// unsubscribe all bidask
	err = sinopacapi.Get().UnSubscribeAllByType(sinopacapi.TickTypeStockBidAsk)
	if err != nil {
		log.Get().Panic(err)
	}

	// sub event targets
	eventbus.Get().SubscribeSubscribeTargets(targetsBusCallback)
	// sub event to unsubscribe
	eventbus.Get().SubscribeUnSubscribeTargets(unSubTargetsBusCallback)
}

func targetsBusCallback(targetArr []*dbagent.Target) {
	var subStockArr []string
	for _, v := range targetArr {
		subStockArr = append(subStockArr, v.Stock.Number)
	}

	// realtime tick
	err := sinopacapi.Get().SubRealTimeTick(subStockArr)
	if err != nil {
		log.Get().Panic(err)
	}
	// realtime bidask
	err = sinopacapi.Get().SubBidAsk(subStockArr)
	if err != nil {
		log.Get().Panic(err)
	}
}

func unSubTargetsBusCallback(target *dbagent.Target) {
	// realtime tick
	err := sinopacapi.Get().UnSubRealTimeTick([]string{target.Stock.Number})
	if err != nil {
		log.Get().Panic(err)
	}
	// realtime bidask
	err = sinopacapi.Get().UnSubBidAsk([]string{target.Stock.Number})
	if err != nil {
		log.Get().Panic(err)
	}
}
