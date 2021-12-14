// Package subscribe package subscribe
package subscribe

import (
	"trade_agent/pkg/dbagent"
	"trade_agent/pkg/eventbus"
	"trade_agent/pkg/log"
)

// InitSubscribe InitSubscribe
func InitSubscribe() {
	// event bus scribe
	_ = eventbus.Get().Subscribe(eventbus.BusTopicTargets(), targetsBusCallback)
}

func targetsBusCallback(tmp *dbagent.Target) error {
	log.Get().Info(tmp.Stock.Number)
	return nil
}
