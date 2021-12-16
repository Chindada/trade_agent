// Package history package history
package history

import (
	"sync"
	"trade_agent/global"
	"trade_agent/pkg/cache"
	"trade_agent/pkg/dbagent"
	"trade_agent/pkg/eventbus"
	"trade_agent/pkg/log"
	"trade_agent/pkg/modules/tradeday"
	"trade_agent/pkg/sinopacapi"
)

// InitHistory InitHistory
func InitHistory() {
	err := eventbus.Get().Sub(eventbus.TopicTargets(), targetsBusCallback)
	if err != nil {
		log.Get().Panic(err)
	}
}

func targetsBusCallback(tmp []*dbagent.Target) error {
	fetchDate := tradeday.GetLastNTradeDayByDate(1, cache.GetCache().GetTradeDay())
	var stockNumArr []string
	for _, v := range tmp {
		stockNumArr = append(stockNumArr, v.Stock.Number)
	}
	subStockClose(stockNumArr, fetchDate[0])

	errChan := make(chan error)
	var wg sync.WaitGroup
	for _, v := range tmp {
		wg.Add(1)
		target := *v
		go func(v dbagent.Target, wg *sync.WaitGroup) {
			defer wg.Done()
			err := sinopacapi.Get().FetchEntireTickByStockAndDate(v.Stock.Number, tradeday.GetLastNTradeDayByDate(1, v.TradeDay)[0].Format(global.ShortTimeLayout))
			if err != nil {
				errChan <- err
			}
		}(target, &wg)
	}
	wg.Wait()
	close(errChan)
	for {
		err, ok := <-errChan
		if err != nil {
			log.Get().Error(err)
		}
		if !ok {
			break
		}
	}
	return nil
}
