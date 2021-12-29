// Package order package order
package order

import (
	"sync"
	"time"
	"trade_agent/pkg/cache"
	"trade_agent/pkg/config"
	"trade_agent/pkg/dbagent"
	"trade_agent/pkg/eventbus"
	"trade_agent/pkg/log"
	"trade_agent/pkg/sinopacapi"
)

var orderLock sync.Mutex

// InitOrder InitOrder
func InitOrder() {
	log.Get().Info("Initial Order")

	err := subOrderStatus()
	if err != nil {
		log.Get().Panic(err)
	}

	err = eventbus.Get().Sub(eventbus.TopicStockOrder(), orderCallback)
	if err != nil {
		log.Get().Panic(err)
	}
}

func orderCallback(order *sinopacapi.Order) error {
	defer orderLock.Unlock()
	orderLock.Lock()
	// check waiting order
	if waitingOrder := cache.GetCache().GetOrderWaiting(order.StockNum); waitingOrder != nil {
		return nil
	}

	switch order.Action {
	case sinopacapi.ActionBuy:
		historyOrderBuy := cache.GetCache().GetOrderBuy(order.StockNum)
		historyOrderSell := cache.GetCache().GetOrderSell(order.StockNum)
		if len(historyOrderBuy) > len(historyOrderSell) {
			return nil
		}
	case sinopacapi.ActionSell:
		historyOrderBuy := cache.GetCache().GetOrderBuy(order.StockNum)
		historyOrderSell := cache.GetCache().GetOrderSell(order.StockNum)
		if len(historyOrderBuy) <= len(historyOrderSell) {
			return nil
		}
	case sinopacapi.ActionSellFirst:
		historyOrderSellFirst := cache.GetCache().GetOrderSellFirst(order.StockNum)
		historyOrderBuyLater := cache.GetCache().GetOrderBuyLater(order.StockNum)
		if len(historyOrderSellFirst) > len(historyOrderBuyLater) {
			return nil
		}
	case sinopacapi.ActionBuyLater:
		historyOrderSellFirst := cache.GetCache().GetOrderSellFirst(order.StockNum)
		historyOrderBuyLater := cache.GetCache().GetOrderBuyLater(order.StockNum)
		if len(historyOrderSellFirst) <= len(historyOrderBuyLater) {
			return nil
		}
	}

	// decide quantiy by history data
	order.Quantity = 1
	orderRes, err := sinopacapi.Get().PlaceOrder(*order)
	if err != nil {
		return err
	}

	if orderID := orderRes.OrderID; orderID != "" {
		order.TradeTime = time.Now()
		order.OrderID = orderID
		cache.GetCache().Set(cache.KeyOrderWaiting(order.StockNum), order)

		// gorutine for check waiting order status
		go checkWaitingOrder(order)
	}
	return nil
}

func checkWaitingOrder(order *sinopacapi.Order) {
	var waitTime int64
	conf, err := config.Get()
	if err != nil {
		log.Get().Panic(err)
	}

	if order.Action == sinopacapi.ActionBuy || order.Action == sinopacapi.ActionSellFirst {
		waitTime = conf.Trade.TradeInWaitTime
	} else {
		waitTime = conf.Trade.TradeOutWaitTime
	}

	for {
		if order.TradeTime.Add(time.Duration(waitTime) * time.Second).Before(time.Now()) {
			break
		}
		time.Sleep(time.Second)
	}

	if status, err := dbagent.Get().GetOrderStatusByOrderID(order.OrderID); err != nil {
		log.Get().Panic(err)
	} else if status != 4 && status != 5 && status != 6 {
		err = sinopacapi.Get().CancelOrder(order.OrderID)
		if err != nil {
			log.Get().Panic(err)
		}
	}
}
