// Package history package history
package history

import (
	"time"
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

	// fill biasrate cache
	err = calculateBiasRate(targetArr, cache.GetCache().GetHistroyCloseRange())
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

func calculateBiasRate(targetArr []*dbagent.Target, fetchDate []time.Time) error {
	for _, stock := range targetArr {
		var closeArr []float64
		for _, date := range fetchDate {
			close := cache.GetCache().GetHistoryClose(stock.Stock.Number, date)
			if close == 0 {
				continue
			}
			closeArr = append(closeArr, close)
		}

		if len(closeArr) != len(fetchDate) {
			log.Get().WithFields(map[string]interface{}{
				"Stock": stock.Stock.Number,
			}).Error("BiasRate Fail")
			eventbus.Get().Pub(eventbus.TopicUnSubscribeTargets(), stock)
			continue
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
