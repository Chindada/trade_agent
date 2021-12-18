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
	fetchDate := cache.GetCache().GetHistroyRange()
	subStockClose(targetArr, fetchDate)

	errChan := make(chan error)
	var w sync.WaitGroup
	for _, v := range targetArr {
		for _, date := range fetchDate {
			exist, err := dbagent.Get().CheckHistoryTickExistByStockNum(date)
			if err != nil {
				return err
			} else if exist {
				log.Get().Infof("%s %s HistoryTick Already Exist", v.Stock.Number, date.Format(global.ShortTimeLayout))
				continue
			}

			w.Add(1)
			stock := v.Stock.Number
			fetchDate := date
			go func(wg *sync.WaitGroup) {
				defer wg.Done()
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
