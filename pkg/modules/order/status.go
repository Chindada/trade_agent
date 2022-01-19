// Package order package order
package order

import (
	"sync"
	"time"
	"trade_agent/pkg/cache"
	"trade_agent/pkg/dbagent"
	"trade_agent/pkg/log"
	"trade_agent/pkg/mqhandler"
	"trade_agent/pkg/pb"
	"trade_agent/pkg/sinopacapi"
)

var mu sync.Mutex

func subOrderStatus() error {
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
		// on start up update once, then every minute update
		updateTradeBalance()

		for range time.Tick(time.Minute) {
			updateTradeBalance()
		}
	}()

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
	defer mu.Unlock()
	mu.Lock()
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
				// order fail or cancel, remove from waiting cache
				cache.GetCache().SetOrderWaiting(v.GetCode(), nil)
				// quota back
				sinopacapi.Get().SetOrderToQuota(*waitingOrder, false)
				displayOrderResult(v.ToOrderStatus())

			case 6:
				// order filled, remove from waiting cache
				cache.GetCache().SetOrderWaiting(v.GetCode(), nil)

				// order filled, add to filled cache by action
				switch waitingOrder.Action {
				case sinopacapi.ActionBuy:
					cache.GetCache().AppendOrderBuy(waitingOrder)
					cache.GetCache().AppendOrderForward(waitingOrder)
				case sinopacapi.ActionSell:
					cache.GetCache().AppendOrderSell(waitingOrder)
					cache.GetCache().AppendOrderForward(waitingOrder)
				case sinopacapi.ActionSellFirst:
					cache.GetCache().AppendOrderSellFirst(waitingOrder)
					cache.GetCache().AppendOrderReverse(waitingOrder)
				case sinopacapi.ActionBuyLater:
					cache.GetCache().AppendOrderBuyLater(waitingOrder)
					cache.GetCache().AppendOrderReverse(waitingOrder)
				}
				displayOrderResult(v.ToOrderStatus())
			}
		}
		saveStatus = append(saveStatus, v.ToOrderStatus())
	}

	err = dbagent.Get().InsertOrUpdateMultiOrderStatus(saveStatus)
	if err != nil {
		log.Get().Error(err)
	}
}

func updateTradeBalance() {
	dbOrderStatus, err := dbagent.Get().GetOrderStatusByDate(cache.GetCache().GetTradeDay())
	if err != nil {
		log.Get().Panic(err)
	}

	tmp := &dbagent.Balance{
		TradeDay: cache.GetCache().GetTradeDay(),
	}

	if len(dbOrderStatus) == 0 {
		return
	}

	for _, order := range dbOrderStatus {
		if order.Status == 6 {
			switch order.Action {
			case 1:
				tmp.OriginalBalance -= sinopacapi.GetStockBuyCost(order.Price, order.Quantity)
			case 2:
				tmp.OriginalBalance += sinopacapi.GetStockSellCost(order.Price, order.Quantity)
			}
			tmp.Discount += sinopacapi.GetStockTradeFeeDiscount(order.Price, order.Quantity)
		}
	}
	tmp.Total = tmp.OriginalBalance + tmp.Discount

	err = dbagent.Get().InsertOrUpdateBalance(tmp)
	if err != nil {
		log.Get().Panic(err)
	}
}
