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
	// save stock close in date range
	err := subStockClose(targetArr, cache.GetCache().GetHistroyCloseRange())
	if err != nil {
		return err
	}

	// save stock tick in date range
	err = subHistoryTick(targetArr, cache.GetCache().GetHistroyTickRange())
	if err != nil {
		return err
	}

	// save stock kbar in date range
	err = subHistoryKbar(targetArr, cache.GetCache().GetHistroyKbarRange())
	if err != nil {
		return err
	}
	return nil
}
