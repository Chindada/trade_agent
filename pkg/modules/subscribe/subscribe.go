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
	// unsubscribe all
	// realtime tick
	err := sinopacapi.Get().UnSubscribeAllByType(sinopacapi.TickTypeStockRealTime)
	if err != nil {
		log.Get().Panic(err)
	}

	// bidask
	err = sinopacapi.Get().UnSubscribeAllByType(sinopacapi.TickTypeStockBidAsk)
	if err != nil {
		log.Get().Panic(err)
	}

	// sub event targets
	err = eventbus.Get().Sub(eventbus.TopicTargets(), targetsBusCallback)
	if err != nil {
		log.Get().Panic(err)
	}
	log.Get().Info("Initial Subscribe")
}

func targetsBusCallback(tmp []*dbagent.Target) error {
	var subStockArr []string
	for _, v := range tmp {
		subStockArr = append(subStockArr, v.Stock.Number)
	}

	// realtime
	err := sinopacapi.Get().SubRealTimeTick(subStockArr)
	if err != nil {
		log.Get().Panic(err)
	}
	return nil
}
