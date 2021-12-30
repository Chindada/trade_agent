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
	// check by config switch
	if !checkSwitch(order) {
		return nil
	}

	// check waiting order
	if waitingOrder := cache.GetCache().GetOrderWaiting(order.StockNum); waitingOrder != nil {
		return nil
	} else if !isGoodPoint(order) {
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
	if quantity := getQuantityByBiasRate(order); quantity != 0 {
		order.Quantity = quantity
	} else {
		return nil
	}

	orderRes, err := sinopacapi.Get().PlaceOrder(*order)
	if err != nil {
		return err
	} else if orderID := orderRes.OrderID; orderID != "" {
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
	if order.Action == sinopacapi.ActionBuy || order.Action == sinopacapi.ActionSellFirst {
		waitTime = config.GetTradeConfig().TradeInWaitTime
	} else {
		waitTime = config.GetTradeConfig().TradeOutWaitTime
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
		for {
			if checkCancelStatus(order.OrderID) {
				break
			}
			time.Sleep(time.Second)
		}
	}
}

func checkCancelStatus(orderID string) bool {
	status, err := dbagent.Get().GetOrderStatusByOrderID(orderID)
	if err != nil {
		log.Get().Panic(err)
	}
	if status == 5 {
		return true
	}
	return false
}

func checkSwitch(order *sinopacapi.Order) bool {
	tradeSwitch := config.GetSwitchConfig()
	switch order.Action {
	case sinopacapi.ActionBuy:
		// get forward remaining orders
		forwardRemaining, total := cache.GetCache().GetOrderForwardCountDetail()
		if tradeSwitch.EnableBuy && forwardRemaining < tradeSwitch.MeanTimeForward && total < tradeSwitch.ForwardMax {
			return true
		}
	case sinopacapi.ActionSell:
		if tradeSwitch.EnableSell {
			return true
		}
	case sinopacapi.ActionSellFirst:
		// get reverse remaining orders
		reverseRemaining, total := cache.GetCache().GetOrderReverseCountDetail()
		if tradeSwitch.EnableSellFirst && reverseRemaining < tradeSwitch.MeanTimeReverse && total < tradeSwitch.ReverseMax {
			return true
		}
	case sinopacapi.ActionBuyLater:
		if tradeSwitch.EnableBuyLater {
			return true
		}
	}
	return false
}

func getQuantityByBiasRate(order *sinopacapi.Order) int64 {
	switch order.Action {
	case sinopacapi.ActionBuy, sinopacapi.ActionSell:
		biasRate := cache.GetCache().GetBiasRate(order.StockNum)
		if biasRate > 4 {
			return 2
		}
		return 1
	case sinopacapi.ActionSellFirst, sinopacapi.ActionBuyLater:
		biasRate := cache.GetCache().GetBiasRate(order.StockNum)
		if biasRate < -4 {
			return 2
		}
		return 1
	}
	return 0
}

func isGoodPoint(order *sinopacapi.Order) bool {
	// historyTickStatus := cache.GetCache().GetStockHistoryTickAnalyze(order.StockNum)
	// historyKbarStatus := cache.GetCache().GetStockHistoryKbarAnalyze(order.StockNum)
	return false
}
