// Package history package history
package history

import (
	"time"

	"trade_agent/global"
	"trade_agent/pkg/cache"
	"trade_agent/pkg/config"
	"trade_agent/pkg/dbagent"
	"trade_agent/pkg/eventbus"
	"trade_agent/pkg/log"
	"trade_agent/pkg/utils"
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

	// fill biasrate cache
	err = calculateQuaterMA(targetArr, cache.GetCache().GetHistroyCloseRange())
	if err != nil {
		log.Get().Error(err)
		return
	}

	// fill biasrate cache
	err = calculateBiasRate(targetArr, cache.GetCache().GetHistroyCloseRange())
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

	for _, v := range targetArr {
		historyTickAnalyzeArr := cache.GetCache().GetStockHistoryTickAnalyze(v.Stock.Number)
		log.Get().WithFields(map[string]interface{}{
			"Stock":  v.Stock.Number,
			"Length": len(historyTickAnalyzeArr),
		}).Info("HistoryTickAnalyzeArr")
	}
}

func calculateBiasRate(targetArr []*dbagent.Target, fetchDate []time.Time) error {
	for _, stock := range targetArr {
		var closeArr []float64
		for _, date := range fetchDate {
			close := cache.GetCache().GetHistoryClose(stock.Stock.Number, date)
			if close == 0 {
				log.Get().Warnf("%s on %s close is 0", stock.Stock.Number, date.Format(global.ShortTimeLayout))
				continue
			}
			closeArr = append(closeArr, close)
		}

		if len(closeArr) != len(fetchDate) {
			log.Get().WithFields(map[string]interface{}{
				"Stock": stock.Stock.Number,
			}).Error("BiasRate Fail")
			eventbus.Get().PublishUnSubscribeTargets(stock)
			continue
		}

		biasRate, err := utils.GetBiasRateByCloseArr(closeArr)
		if err != nil {
			return err
		}

		cache.GetCache().SetBiasRate(stock.Stock.Number, biasRate)
		log.Get().WithFields(map[string]interface{}{
			"Stock": stock.Stock.Number,
			"Value": biasRate,
		}).Info("BiasRate")
	}
	return nil
}

func calculateQuaterMA(targetArr []*dbagent.Target, fetchDate []time.Time) error {
	analyzeConf := config.GetAnalyzeConfig()
	var closeArr []float64
	for _, stock := range targetArr {
		closeArr = []float64{}
		for _, date := range fetchDate {
			close := cache.GetCache().GetHistoryClose(stock.Stock.Number, date)
			if close == 0 {
				log.Get().Warnf("%s on %s close is 0", stock.Stock.Number, date.Format(global.ShortTimeLayout))
				break
			}
			closeArr = append(closeArr, close)
		}
		i := 0
		for {
			if i+int(analyzeConf.MAPeriod) > len(closeArr) {
				break
			}
			tmp := closeArr[i : i+int(analyzeConf.MAPeriod)]
			ma, err := utils.GenerareMAByCloseArr(tmp)
			if err != nil {
				return err
			}
			calendarDate, err := dbagent.Get().GetCalendarDate(fetchDate[i])
			if err != nil {
				log.Get().Panic(err)
			}
			toDBMA := &dbagent.HistoryMA{
				QuaterMA:     utils.Round(ma, 2),
				Stock:        stock.Stock,
				CalendarDate: calendarDate,
			}

			if err := dbagent.Get().InsertOrUpdateHistoryMA(toDBMA); err != nil {
				log.Get().Panic(err)
			}
			i++
		}
	}
	return nil
}
