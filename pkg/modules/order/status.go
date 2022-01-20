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
			saveStatus = append(saveStatus, v.ToOrderStatus())
		}
	}

	if len(saveStatus) == 0 {
		return
	}

	err = dbagent.Get().InsertOrUpdateMultiOrderStatus(saveStatus)
	if err != nil {
		log.Get().Error(err)
	}
}

func updateTradeBalance() {
	forwardOrder := cache.GetCache().GetOrderForward()
	reverseOrder := cache.GetCache().GetOrderReverse()

	if len(forwardOrder)/2 == 0 && len(reverseOrder)/2 == 0 {
		return
	}

	var forwardBalance, revereBalance, discount int64
	for _, v := range forwardOrder {
		switch v.Action {
		case sinopacapi.ActionBuy:
			forwardBalance -= sinopacapi.GetStockBuyCost(v.Price, v.Quantity)
		case sinopacapi.ActionSell:
			forwardBalance += sinopacapi.GetStockSellCost(v.Price, v.Quantity)
		}
		discount += sinopacapi.GetStockTradeFeeDiscount(v.Price, v.Quantity)
	}

	for _, v := range reverseOrder {
		switch v.Action {
		case sinopacapi.ActionSellFirst:
			revereBalance += sinopacapi.GetStockSellCost(v.Price, v.Quantity)
		case sinopacapi.ActionBuyLater:
			revereBalance -= sinopacapi.GetStockBuyCost(v.Price, v.Quantity)
		}
		discount += sinopacapi.GetStockTradeFeeDiscount(v.Price, v.Quantity)
	}

	tmp := &dbagent.Balance{
		TradeDay:        cache.GetCache().GetTradeDay(),
		TradeCount:      int64(len(forwardOrder)/2 + len(reverseOrder)/2),
		Forward:         forwardBalance,
		Reverse:         revereBalance,
		OriginalBalance: forwardBalance + revereBalance,
		Discount:        discount,
		Total:           forwardBalance + revereBalance + discount,
	}

	err := dbagent.Get().InsertOrUpdateBalance(tmp)
	if err != nil {
		log.Get().Panic(err)
	}
}
