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

	"google.golang.org/protobuf/proto"
)

var wg sync.WaitGroup

// InitStock InitStock
func InitStock() {
	err := subStockDeatail()
	if err != nil {
		log.Get().Panic(err)
	}

	// add stock detail to cache from db
	inDBStock, err := dbagent.Get().GetAllStockMap()
	if err != nil {
		log.Get().Panic(err)
	}

	// save stock detail to cahce
	for key := range inDBStock {
		cache.GetCache().Set(cache.KeyStockDetail(key), inDBStock[key])
	}
	log.Get().Info("Initial Stock")
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
	if err = proto.Unmarshal(m.Payload(), &body); err != nil {
		log.Get().Errorf("Format Wrong: %s", string(m.Payload()))
		return
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
			continue
		}
		saveStock = append(saveStock, v.ToStock())
		insert++
	}
	// insert
	err = dbagent.Get().InsertMultiStock(saveStock)
	if err != nil {
		log.Get().Panic(err)
	}
	log.Get().WithFields(map[string]interface{}{
		"Exist":  exist,
		"Insert": insert,
	}).Info("GetAllStockDetail")
}
