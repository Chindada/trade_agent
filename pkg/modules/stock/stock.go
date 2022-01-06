// Package stock package stock
package stock

import (
	"sync"
	"trade_agent/pkg/cache"
	"trade_agent/pkg/dbagent"
	"trade_agent/pkg/log"
	"trade_agent/pkg/mqhandler"
	"trade_agent/pkg/pb"
	"trade_agent/pkg/sinopacapi"
)

var wg sync.WaitGroup

// InitStock InitStock
func InitStock() {
	log.Get().Info("Initial Stock")

	err := subStockDeatail()
	if err != nil {
		log.Get().Panic(err)
	}

	// add stock detail to cache from db
	inDBStock, err := dbagent.Get().GetAllDayTradeStockMap()
	if err != nil {
		log.Get().Panic(err)
	}

	// save stock detail to cahce
	for key := range inDBStock {
		cache.GetCache().Set(cache.KeyStockDetail(key), inDBStock[key])
	}
}

func subStockDeatail() error {
	handler := mqhandler.Get()
	wg.Add(1)
	err := handler.Sub(mqhandler.MQSubBody{
		MQTopic:  mqhandler.TopicStockDetail(),
		Once:     true,
		Callback: stockDetailCallback,
	})
	if err != nil {
		return err
	}
	err = sinopacapi.Get().FetchAllStockDetail()
	if err != nil {
		return err
	}
	// wait stock callback return
	wg.Wait()
	return nil
}

// process mq back stock deail, check db record to decide to insert, and add to cache
func stockDetailCallback(m mqhandler.MQMessage) {
	defer wg.Done()
	var err error
	body := pb.StockDetailResponse{}
	err = body.UnmarshalProto(m.Payload())
	if err != nil {
		log.Get().Panic(err)
	}

	var inDBStock map[string]*dbagent.Stock
	inDBStock, err = dbagent.Get().GetAllStockMap()
	if err != nil {
		log.Get().Panic(err)
	}

	var saveStock []*dbagent.Stock
	var exist, insert int
	for _, v := range body.GetStock() {
		// check whether already in db
		if _, ok := inDBStock[v.GetCode()]; ok {
			exist++
			saveStock = append(saveStock, inDBStock[v.GetCode()])
		} else {
			insert++
			saveStock = append(saveStock, v.ToStock())
		}
	}

	// make sure every time startup all stock in db will be day trade
	err = dbagent.Get().UpdateAllStockDayTradeFalse()
	if err != nil {
		log.Get().Panic(err)
	}

	// insert
	err = dbagent.Get().InsertOrUpdateMultiStock(saveStock)
	if err != nil {
		log.Get().Panic(err)
	}

	log.Get().WithFields(map[string]interface{}{
		"Exist":  exist,
		"Insert": insert,
	}).Info("GetAllStockDetail")
}
