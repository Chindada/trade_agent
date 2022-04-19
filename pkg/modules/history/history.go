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
	log.Get().Info("Initial History")

	eventbus.Get().SubscribeTargets(targetsBusCallback)
}

func targetsBusCallback(targetArr []*dbagent.Target) {
	// save stock kbar in date range
	err := subHistoryKbar(targetArr, cache.GetCache().GetHistroyKbarRange())
	if err != nil {
		log.Get().Error(err)
		return
	}

	// save stock close in date range
	err = subStockClose(targetArr, cache.GetCache().GetHistroyCloseRange())
	if err != nil {
		log.Get().Error(err)
		return
	}

	// save stock tick in date range
	err = subHistoryTick(targetArr, cache.GetCache().GetHistroyTickRange())
	if err != nil {
		log.Get().Error(err)
		return
	}

	// send to analyze
	eventbus.Get().PublishNeedAnalyzeTargets(targetArr)

	for _, v := range targetArr {
		historyTickAnalyzeArr := cache.GetCache().GetStockHistoryTickAnalyze(v.Stock.Number)
		log.Get().WithFields(map[string]interface{}{
			"Stock":  v.Stock.Number,
			"Length": len(historyTickAnalyzeArr),
		}).Info("HistoryTickAnalyzeArr")
	}
}
