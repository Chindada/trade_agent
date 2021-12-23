// Package history package history
package history

import (
	"trade_agent/pkg/cache"
	"trade_agent/pkg/dbagent"
	"trade_agent/pkg/eventbus"
	"trade_agent/pkg/log"
)

// InitHistory InitHistory
func InitHistory() {
	err := eventbus.Get().Sub(eventbus.TopicTargets(), targetsBusCallback)
	if err != nil {
		log.Get().Panic(err)
	}
	log.Get().Info("Initial History")
}

func targetsBusCallback(targetArr []*dbagent.Target) error {
	// get date range for fetch
	fetchDate := cache.GetCache().GetHistroyRange()

	// update stock close in date range
	err := subStockClose(targetArr, fetchDate)
	if err != nil {
		return err
	}
	err = subHistoryTick(targetArr, fetchDate)
	if err != nil {
		return err
	}
	err = subHistoryKbar(targetArr, fetchDate)
	if err != nil {
		return err
	}
	return nil
}
