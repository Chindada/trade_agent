// Package analyze package analyze
package analyze

import (
	"sync"
	"time"

	"trade_agent/global"
	"trade_agent/pkg/cache"
	"trade_agent/pkg/config"
	"trade_agent/pkg/dbagent"
	"trade_agent/pkg/eventbus"
	"trade_agent/pkg/log"
	"trade_agent/pkg/modules/tradeday"
	"trade_agent/pkg/utils"
)

// InitAnalyze InitAnalyze
func InitAnalyze() {
	log.Get().Info("Initial Analyze")

	belowQuaterMap = make(map[time.Time][]dbagent.Stock)
	lastBelowMAStock = make(map[string]*dbagent.HistoryMA)

	eventbus.Get().SubscribeNeedAnalyzeTargets(targetsBusCallback)
}

func targetsBusCallback(targetArr []*dbagent.Target) {
	// fill biasrate cache
	err := calculateQuaterMA(targetArr, cache.GetCache().GetHistroyCloseRange())
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

	err = findBelowQuaterMATargets(targetArr)
	if err != nil {
		log.Get().Error(err)
		return
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
			calendarDate := cache.GetCache().GetCalendarID(fetchDate[i].Format(global.ShortTimeLayout))
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

var (
	lastBelowMAStock map[string]*dbagent.HistoryMA
	belowQuaterMap   map[time.Time][]dbagent.Stock
	belowQuaterLock  sync.Mutex
)

// GetBelowQuaterMap GetBelowQuaterMap
func GetBelowQuaterMap() map[time.Time][]dbagent.Stock {
	belowQuaterLock.Lock()
	tmp := belowQuaterMap
	if len(lastBelowMAStock) != 0 {
		for _, s := range lastBelowMAStock {
			if open := cache.GetCache().GetHistoryOpen(s.Stock.Number, cache.GetCache().GetTradeDay()); open != 0 {
				if open > s.QuaterMA {
					belowQuaterMap[s.CalendarDate.Date] = append(belowQuaterMap[s.CalendarDate.Date], *s.Stock)
				}
				delete(lastBelowMAStock, s.Stock.Number)
			}
		}
	}
	belowQuaterLock.Unlock()
	return tmp
}

func findBelowQuaterMATargets(targetArr []*dbagent.Target) error {
	defer log.Get().Info("FindBelowQuaterMATargets Done")
	defer belowQuaterLock.Unlock()
	belowQuaterLock.Lock()
	for _, t := range targetArr {
		tmp := *t
		maArr, err := dbagent.Get().GetAllQuaterMAByStockID(int64(tmp.Stock.ID))
		if err != nil {
			log.Get().Error(err)
		}

		for _, ma := range maArr {
			if close := cache.GetCache().GetHistoryClose(ma.Stock.Number, ma.CalendarDate.Date); close != 0 && close-ma.QuaterMA < 0 {
				nextTradeDay, err := tradeday.GetAbsNextTradeDayTime(ma.CalendarDate.Date)
				if err != nil {
					log.Get().Error(err)
				}
				if nextTradeDay.Equal(cache.GetCache().GetTradeDay()) {
					tmp := ma
					lastBelowMAStock[tmp.Stock.Number] = &tmp
				}
				if nextOpen := cache.GetCache().GetHistoryOpen(ma.Stock.Number, nextTradeDay); nextOpen != 0 && nextOpen-ma.QuaterMA > 0 {
					belowQuaterMap[ma.CalendarDate.Date] = append(belowQuaterMap[ma.CalendarDate.Date], *tmp.Stock)
				}
			}
		}
	}
	return nil
}
