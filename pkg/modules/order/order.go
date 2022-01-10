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

	// find all target's available action, and send order
	// price come from bidask best price
	go clearAllUnFinished()

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
		order.OrderID = orderID
		cache.GetCache().SetOrderWaiting(order.StockNum, order)

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
	isOpen := cache.GetCache().GetIsOpenWithEndWaitTime()
	switch order.Action {
	case sinopacapi.ActionBuy:
		// get forward remaining orders
		forwardRemaining, total := cache.GetCache().GetOrderForwardCountDetail()
		if tradeSwitch.EnableBuy && forwardRemaining < tradeSwitch.MeanTimeForward && total < tradeSwitch.ForwardMax && isOpen {
			return true
		}
	case sinopacapi.ActionSell:
		if tradeSwitch.EnableSell {
			return true
		}
	case sinopacapi.ActionSellFirst:
		// get reverse remaining orders
		reverseRemaining, total := cache.GetCache().GetOrderReverseCountDetail()
		if tradeSwitch.EnableSellFirst && reverseRemaining < tradeSwitch.MeanTimeReverse && total < tradeSwitch.ReverseMax && isOpen {
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
	// historyKbarStatus := cache.GetCache().GetStockHistoryKbarAnalyze(order.StockNum)
	return true
}

// clearAllUnFinished clearAllUnFinished
func clearAllUnFinished() {
	tradeOutEndTime := cache.GetCache().GetTradeDayTradeOutEndTime()
	for {
		if time.Now().Before(tradeOutEndTime) {
			continue
		}

		targetArr := cache.GetCache().GetTargets()
		for _, t := range targetArr {
			historyOrderBuy := cache.GetCache().GetOrderBuy(t.Stock.Number)
			historyOrderSell := cache.GetCache().GetOrderSell(t.Stock.Number)
			if len(historyOrderBuy) > len(historyOrderSell) {
				order := &sinopacapi.Order{
					StockNum:  t.Stock.Number,
					Price:     cache.GetCache().GetRealTimeTickClose(t.Stock.Number),
					Action:    sinopacapi.ActionSell,
					TradeTime: time.Now(),
				}
				eventbus.Get().Pub(eventbus.TopicStockOrder(), order)
			}

			historyOrderSellFirst := cache.GetCache().GetOrderSellFirst(t.Stock.Number)
			historyOrderBuyLater := cache.GetCache().GetOrderBuyLater(t.Stock.Number)
			if len(historyOrderSellFirst) > len(historyOrderBuyLater) {
				order := &sinopacapi.Order{
					StockNum:  t.Stock.Number,
					Price:     cache.GetCache().GetRealTimeTickClose(t.Stock.Number),
					Action:    sinopacapi.ActionBuyLater,
					TradeTime: time.Now(),
				}
				eventbus.Get().Pub(eventbus.TopicStockOrder(), order)
			}
		}
		time.Sleep(15 * time.Second)
	}
}
