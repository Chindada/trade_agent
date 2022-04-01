// Package history package history
package history

import (
	"sync"
	"time"

	"trade_agent/global"
	"trade_agent/pkg/cache"
	"trade_agent/pkg/config"
	"trade_agent/pkg/dbagent"
	"trade_agent/pkg/log"
	"trade_agent/pkg/sinopacapi"
)

func subHistoryTick(targetArr []*dbagent.Target, fetchDate []time.Time) error {
	// check history tick exist or fetch
	errChan := make(chan error)
	var w sync.WaitGroup
	for _, v := range targetArr {
		for _, date := range fetchDate {
			if exist, dbErr := dbagent.Get().CheckHistoryTickExistByStockID(v.StockID, date); dbErr != nil {
				return dbErr
			} else if exist {
				log.Get().WithFields(map[string]interface{}{
					"Stock": v.Stock.Number,
					"Date":  date.Format(global.ShortTimeLayout),
				}).Info("HistoryTick Already Exist")

				// select from db to analyze to cache
				analyzeConf := config.GetAnalyzeConfig()
				if dbHistoryTick, fetchErr := dbagent.Get().GetHistoryTickByStockIDAndDate(v.StockID, date); fetchErr != nil {
					return fetchErr
				} else if analyzeVolumeArr := dbHistoryTick.Analyzer(analyzeConf.TickAnalyzeMinPeriod, analyzeConf.TickAnalyzeMaxPeriod); len(analyzeVolumeArr) != 0 {
					cache.GetCache().AppendHistoryTickAnalyze(v.Stock.Number, analyzeVolumeArr)
				}
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
				if sinoErr := sinopacapi.Get().FetchHistoryTickByStockAndDate(stock, fetchDate.Format(global.ShortTimeLayout)); sinoErr != nil {
					errChan <- sinoErr
				}
			}(&w)
		}
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
