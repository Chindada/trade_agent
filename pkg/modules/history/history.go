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

	err := eventbus.Get().Sub(eventbus.TopicTargets(), targetsBusCallback)
	if err != nil {
		log.Get().Panic(err)
	}
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

	// fill biasrate cache
	for _, stock := range targetArr {
		var closeArr []float64
		for _, date := range cache.GetCache().GetHistroyCloseRange() {
			close := cache.GetCache().GetHistoryClose(stock.Stock.Number, date)
			closeArr = append(closeArr, close)
		}
		biasRate, err := getBiasRateByCloseArr(closeArr)
		if err != nil {
			return err
		}
		cache.GetCache().Set(cache.KeyBiasRate(stock.Stock.Number), biasRate)
		log.Get().WithFields(map[string]interface{}{
			"Stock": stock.Stock.Number,
			"Value": biasRate,
		}).Info("BiasRate")
	}

	return nil
}
