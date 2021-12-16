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
	err := sinopacapi.Get().UnSubscribeAllByType(sinopacapi.StreamType)
	if err != nil {
		log.Get().Panic(err)
	}
	err = eventbus.Get().Sub(eventbus.TopicTargets(), targetsBusCallback)
	if err != nil {
		log.Get().Panic(err)
	}
}

func targetsBusCallback(tmp []*dbagent.Target) error {
	var subStockArr []string
	for _, v := range tmp {
		subStockArr = append(subStockArr, v.Stock.Number)
	}
	err := sinopacapi.Get().SubStreamTick(subStockArr)
	if err != nil {
		log.Get().Panic(err)
	}
	return nil
}
