// Package history package history
package history

import (
	"sync"
	"trade_agent/global"
	"trade_agent/pkg/cache"
	"trade_agent/pkg/dbagent"
	"trade_agent/pkg/eventbus"
	"trade_agent/pkg/log"
	"trade_agent/pkg/sinopacapi"
)

// InitHistory InitHistory
func InitHistory() {
	err := eventbus.Get().Sub(eventbus.TopicTargets(), targetsBusCallback)
	if err != nil {
		log.Get().Panic(err)
	}
	log.Get().Info("Initial History")
}

func targetsBusCallback(targetArr []*dbagent.Target) error {
	// get date range for fetch
	fetchDate := cache.GetCache().GetHistroyRange()

	// update stock close in date range
	subStockClose(targetArr, fetchDate)

	// check history tick exist or fetch
	errChan := make(chan error)
	var w sync.WaitGroup
	for _, v := range targetArr {
		for _, date := range fetchDate {
			exist, err := dbagent.Get().CheckHistoryTickExistByStockNum(date)
			if err != nil {
				return err
			} else if exist {
				log.Get().WithFields(map[string]interface{}{
					"Stock": v.Stock.Number,
					"Date":  date.Format(global.ShortTimeLayout),
				}).Info("HistoryTick Already Exist")
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
				err := sinopacapi.Get().FetchHistoryTickByStockAndDate(stock, fetchDate.Format(global.ShortTimeLayout))
				if err != nil {
					errChan <- err
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
