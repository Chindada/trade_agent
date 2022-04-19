// Package history package history
package history

import (
	"sync"
	"time"

	"trade_agent/global"
	"trade_agent/pkg/cache"
	"trade_agent/pkg/dbagent"
	"trade_agent/pkg/log"
	"trade_agent/pkg/sinopacapi"
)

type tempFetch struct {
	target *dbagent.Target
	date   time.Time
}

func subHistoryKbar(targetArr []*dbagent.Target, fetchDate []time.Time) error {
	errChan := make(chan error)
	var w sync.WaitGroup
	var stockIDArr []uint
	targetMap := make(map[uint]*dbagent.Target)
	for _, v := range targetArr {
		stockIDArr = append(stockIDArr, v.Stock.ID)
		targetMap[v.Stock.ID] = v
	}
	noDataDateArr := []tempFetch{}
	for _, date := range fetchDate {
		kbarArr, err := dbagent.Get().GetHistoryOpenByMultiStockIDAndDate(stockIDArr, date)
		if err != nil {
			return err
		}
		for s, arr := range kbarArr {
			if len(arr) == 0 {
				noDataDateArr = append(noDataDateArr, tempFetch{targetMap[s], date})
			} else {
				cache.GetCache().SetStockHistoryOpen(arr[0].Stock.Number, arr[0].Open, date)
				var close, open, high, low float64
				var volume int64
				var lastTickTime time.Time
				for i, kbar := range arr {
					if i == 0 {
						open = kbar.Open
					}
					if i == len(arr)-1 {
						close = kbar.Close
						lastTickTime = kbar.TickTime
					}
					if high == 0 {
						high = kbar.High
					} else if kbar.High > high {
						high = kbar.High
					}
					if low == 0 {
						low = kbar.Low
					} else if kbar.Low < low {
						low = kbar.Low
					}
					volume += kbar.Volume
				}
				cache.GetCache().SetStockHistoryDayKbar(date.Format(global.ShortTimeLayout), arr[0].Stock.Number, &dbagent.HistoryKbar{
					Close:    close,
					Open:     open,
					High:     high,
					Low:      low,
					Volume:   volume,
					TickTime: lastTickTime,
				})
				log.Get().WithFields(map[string]interface{}{
					"Stock": arr[0].Stock.Number,
					"Date":  date.Format(global.ShortTimeLayout),
				}).Info("HistoryKbar Already Exist")
			}
		}
	}
	// 	// status := dbHistoryKbar.Analyzer()
	// 	// cache.GetCache().SetStockHistoryKbarAnalyze(v.Stock.Number, status)
	// 	continue
	// }
	// // does not exist, fetch.
	for _, fetch := range noDataDateArr {
		w.Add(1)
		stock := fetch.target.Stock.Number
		fetchDate := fetch.date.Format(global.ShortTimeLayout)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			log.Get().WithFields(map[string]interface{}{
				"Stock": stock,
				"Date":  fetchDate,
			}).Info("Fetching HistoryKbar")
			sinoErr := sinopacapi.Get().FetchHistoryKbarByDateRange(stock, fetchDate, fetchDate)
			if sinoErr != nil {
				errChan <- sinoErr
			}
		}(&w)
	}
	w.Wait()
	close(errChan)
	for {
		err, ok := <-errChan
		if !ok {
			break
		}
		if err != nil {
			log.Get().Error(err)
			return err
		}
	}
	return nil
}
