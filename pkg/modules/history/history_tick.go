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

func subHistoryTick(targetArr []*dbagent.Target, fetchDate []time.Time) error {
	var err error
	// check history tick exist or fetch
	errChan := make(chan error)
	var w sync.WaitGroup
	for _, v := range targetArr {
		for _, date := range fetchDate {
			var exist bool
			exist, err = dbagent.Get().CheckHistoryTickExistByStockNum(date)
			if err != nil {
				return err
			} else if exist {
				log.Get().WithFields(map[string]interface{}{
					"Stock": v.Stock.Number,
					"Date":  date.Format(global.ShortTimeLayout),
				}).Info("HistoryTick Already Exist")

				// select from db to analyze to cache
				var dbHistoryTick dbagent.HistoryTickArr
				dbHistoryTick, err = dbagent.Get().GetHistoryTickByStockIDAndDate(v.StockID, date)
				if err != nil {
					return err
				}
				analyzeVolumeArr := dbHistoryTick.Analyzer()
				cache.GetCache().Set(cache.KeyStockHistoryTickAnalyze(v.Stock.Number), analyzeVolumeArr)
				continue
			}
			// does not exist, fetch.
			w.Add(1)
			stock := v.Stock.Number
			fetchDate := date
			go func(wg *sync.WaitGroup) {
				defer wg.Done()
				log.Get().WithFields(map[string]interface{}{
					"Stock": stock,
					"Date":  fetchDate.Format(global.ShortTimeLayout),
				}).Info("Fetching HistoryTick")
				err = sinopacapi.Get().FetchHistoryTickByStockAndDate(stock, fetchDate.Format(global.ShortTimeLayout))
				if err != nil {
					errChan <- err
				}
			}(&w)
		}
	}
	w.Wait()
	close(errChan)
	for {
		var ok bool
		err, ok = <-errChan
		if !ok {
			break
		}
		if err != nil {
			log.Get().Error(err)
			return err
		}
	}
	return err
}
