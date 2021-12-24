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
	err = eventbus.Get().Sub(eventbus.TopicSubscribeTargets(), targetsBusCallback)
	if err != nil {
		log.Get().Panic(err)
	}
}

func targetsBusCallback(targetArr []*dbagent.Target) error {
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

	return nil
}
