// Package order package order
package order

import (
	"sync"
	"trade_agent/pkg/cache"
	"trade_agent/pkg/eventbus"
	"trade_agent/pkg/log"
	"trade_agent/pkg/sinopacapi"
)

var orderLock sync.Mutex

// InitOrder InitOrder
func InitOrder() {
	log.Get().Info("Initial Order")

	err := updateOrderStatus()
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
	case sinopacapi.ActionSellFirst:
		historyOrderSellFirst := cache.GetCache().GetOrderSellFirst(order.StockNum)
		historyOrderBuyLater := cache.GetCache().GetOrderBuyLater(order.StockNum)
		if len(historyOrderSellFirst) > len(historyOrderBuyLater) {
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
		cache.GetCache().Set(cache.KeyOrderWaiting(order.StockNum), order)
	}
	return nil
}
