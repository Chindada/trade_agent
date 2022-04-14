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
		tickArrMap, err := dbagent.Get().GetHistoryTickByMultiStockIDAndDate(stockIDArr, date)
		if err != nil {
			return err
		}

		for s, arr := range tickArrMap {
			if len(arr) == 0 {
				noDataDateArr = append(noDataDateArr, tempFetch{targetMap[s], date})
			} else {
				analyzeConf := config.GetAnalyzeConfig()
				if analyzeVolumeArr := arr.Analyzer(analyzeConf.TickAnalyzeMinPeriod, analyzeConf.TickAnalyzeMaxPeriod); len(analyzeVolumeArr) != 0 {
					cache.GetCache().AppendHistoryTickAnalyze(arr[0].Stock.Number, analyzeVolumeArr)
				}
				log.Get().WithFields(map[string]interface{}{
					"Stock": arr[0].Stock.Number,
					"Date":  date.Format(global.ShortTimeLayout),
				}).Info("HistoryTick Already Exist")
			}
		}
	}

	for _, fetch := range noDataDateArr {
		w.Add(1)
		stock := fetch.target.Stock.Number
		fetchDate := fetch.date
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
