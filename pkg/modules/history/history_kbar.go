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

func subHistoryKbar(targetArr []*dbagent.Target, fetchDate []time.Time) error {
	var err error
	// check history tick exist or fetch
	errChan := make(chan error)
	var w sync.WaitGroup
	for _, v := range targetArr {
		for _, date := range fetchDate {
			var exist bool
			exist, err = dbagent.Get().CheckHistoryKbarExistByStockNum(date)
			if err != nil {
				return err
			} else if exist {
				log.Get().WithFields(map[string]interface{}{
					"Stock": v.Stock.Number,
					"Date":  date.Format(global.ShortTimeLayout),
				}).Info("HistoryKbar Already Exist")

				// select from db to analyze to cache
				var dbHistoryKbar dbagent.HistoryKbarArr
				dbHistoryKbar, err = dbagent.Get().GetHistoryKbarByStockIDAndDate(v.StockID, date)
				if err != nil {
					return err
				}
				status := dbHistoryKbar.Analyzer()
				cache.GetCache().Set(cache.KeyStockHistoryKbarAnalyze(v.Stock.Number), status)
				continue
			}
			// does not exist, fetch.
			w.Add(1)
			stock := v.Stock.Number
			fetchDate := date.Format(global.ShortTimeLayout)
			go func(wg *sync.WaitGroup) {
				defer wg.Done()
				log.Get().WithFields(map[string]interface{}{
					"Stock": stock,
					"Date":  fetchDate,
				}).Info("Fetching HistoryKbar")
				err = sinopacapi.Get().FetchHistoryKbarByDateRange(stock, fetchDate, fetchDate)
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
