// Package order package order
package order

import (
	"time"
	"trade_agent/pkg/cache"
	"trade_agent/pkg/dbagent"
	"trade_agent/pkg/log"
	"trade_agent/pkg/mqhandler"
	"trade_agent/pkg/pb"
	"trade_agent/pkg/sinopacapi"
)

func updateOrderStatus() error {
	handler := mqhandler.Get()
	err := handler.Sub(mqhandler.MQSubBody{
		MQTopic:  mqhandler.TopicOrderStatus(),
		Once:     false,
		Callback: orderStausCallback,
	})
	if err != nil {
		return err
	}
	go func() {
		for range time.Tick(1*time.Second + 500*time.Millisecond) {
			if err := sinopacapi.Get().FetchOrderStatus(); err != nil {
				log.Get().Error(err)
			}
		}
	}()
	return nil
}

func orderStausCallback(m mqhandler.MQMessage) {
	body := pb.OrderStatusHistoryResponse{}
	err := body.UnmarshalProto(m.Payload())
	if err != nil {
		log.Get().Panic(err)
	}
	var saveStatus []*dbagent.OrderStatus
	for _, v := range body.GetData() {
		// check waiting order
		if waitingOrder := cache.GetCache().GetOrderWaiting(v.GetCode()); waitingOrder != nil && v.GetOrderId() == waitingOrder.OrderID {
			statusMap := dbagent.StatusListMap
			switch statusMap[v.GetStatus()] {
			case 4, 5:
				cache.GetCache().Set(cache.KeyOrderWaiting(v.GetCode()), nil)
			case 6:
				cache.GetCache().Set(cache.KeyOrderWaiting(v.GetCode()), nil)
				var cacheKey string
				switch waitingOrder.Action {
				case sinopacapi.ActionBuy:
					cacheKey = cache.KeyOrderBuy(waitingOrder.StockNum)
				case sinopacapi.ActionSell:
					cacheKey = cache.KeyOrderSell(waitingOrder.StockNum)
				case sinopacapi.ActionSellFirst:
					cacheKey = cache.KeyOrderSellFirst(waitingOrder.StockNum)
				case sinopacapi.ActionBuyLater:
					cacheKey = cache.KeyOrderBuyLater(waitingOrder.StockNum)
				}
				waitingOrder.TradeTime = time.Now()
				cache.GetCache().Set(cacheKey, waitingOrder)
			}
		}
		saveStatus = append(saveStatus, v.ToOrderStatus())
	}

	err = dbagent.Get().InsertOrUpdateMultiOrderStatus(saveStatus)
	if err != nil {
		log.Get().Error(err)
	}
}
