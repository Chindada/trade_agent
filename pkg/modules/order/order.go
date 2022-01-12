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

func orderCallback(order *sinopacapi.Order) {
	defer orderLock.Unlock()
	orderLock.Lock()
	// check by config switch
	if !checkSwitch(order) {
		return
	}

	// check waiting order, if still waiting, return
	if waitingOrder := cache.GetCache().GetOrderWaiting(order.StockNum); waitingOrder != nil {
		return
	} else if !isGoodPoint(order) {
		return
	}

	// decide quantiy by history data
	quantity := getQuantityByBiasRate(order)
	if quantity == 0 {
		return
	}
	order.Quantity = quantity

	orderRes, err := sinopacapi.Get().PlaceOrder(*order)
	if err != nil || orderRes.Status != sinopacapi.StatusSuccuss {
		if err != nil {
			log.Get().Error(err)
		}
		log.Get().WithFields(map[string]interface{}{
			"Stock":  order.StockNum,
			"Action": order.Action,
		}).Error("PlaceOrder Fail")
		return
	}

	if orderID := orderRes.OrderID; orderID != "" {
		order.OrderID = orderID
		cache.GetCache().SetOrderWaiting(order.StockNum, order)

		// gorutine for check waiting order status
		go checkWaitingOrder(order)
	}
}

func checkWaitingOrder(order *sinopacapi.Order) {
	tradeConf := config.GetTradeConfig()

	var waitTime time.Duration
	switch order.Action {
	case sinopacapi.ActionBuy, sinopacapi.ActionSellFirst:
		waitTime = time.Duration(tradeConf.TradeInWaitTime) * time.Second
	case sinopacapi.ActionSell, sinopacapi.ActionBuyLater:
		waitTime = time.Duration(tradeConf.TradeOutWaitTime) * time.Second
	}

	for {
		if order.TradeTime.Add(waitTime).Before(time.Now()) {
			break
		}
		time.Sleep(5 * time.Second)
	}

	statusMap := dbagent.StatusListMap
	status, err := sinopacapi.Get().FetchOrderStatusByOrderID(order.OrderID)
	if err != nil {
		log.Get().Panic(err)
	}

	if statusMap[status] != 4 && statusMap[status] != 5 && statusMap[status] != 6 {
		err = sinopacapi.Get().CancelOrder(order.OrderID)
		if err != nil {
			log.Get().Error(err)
			if isOrderNotCanceld(order.OrderID) {
				checkWaitingOrder(order)
			}
		}
	}
}

func isOrderNotCanceld(orderID string) bool {
	statusMap := dbagent.StatusListMap
	status, err := sinopacapi.Get().FetchOrderStatusByOrderID(orderID)
	if err != nil {
		log.Get().Error(err)
		return false
	}
	if statusMap[status] == 4 || statusMap[status] == 5 || statusMap[status] == 6 {
		return false
	}
	return true
}

func checkSwitch(order *sinopacapi.Order) bool {
	switchConf := config.GetSwitchConfig()
	isAllowTrade := cache.GetCache().GetIsAllowTrade()

	switch order.Action {
	case sinopacapi.ActionBuy:
		// get forward remaining orders
		forwardRemaining, total := cache.GetCache().GetOrderForwardCountDetail()
		if switchConf.Buy && forwardRemaining < switchConf.MeanTimeForward && total < switchConf.ForwardMax && isAllowTrade {
			return true
		}
	case sinopacapi.ActionSell:
		if switchConf.Sell {
			return true
		}
	case sinopacapi.ActionSellFirst:
		// get reverse remaining orders
		reverseRemaining, total := cache.GetCache().GetOrderReverseCountDetail()
		if switchConf.SellFirst && reverseRemaining < switchConf.MeanTimeReverse && total < switchConf.ReverseMax && isAllowTrade {
			return true
		}
	case sinopacapi.ActionBuyLater:
		if switchConf.BuyLater {
			return true
		}
	}
	return false
}

func getQuantityByBiasRate(order *sinopacapi.Order) int64 {
	switch order.Action {
	case sinopacapi.ActionBuy, sinopacapi.ActionSell:
		if biasRate := cache.GetCache().GetBiasRate(order.StockNum); biasRate > 4 {
			return 2
		}
		return 1
	case sinopacapi.ActionSellFirst, sinopacapi.ActionBuyLater:
		if biasRate := cache.GetCache().GetBiasRate(order.StockNum); biasRate < -4 {
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
	tradeDayOpenEndTime := cache.GetCache().GetTradeDayOpenEndTime()

	for {
		time.Sleep(15 * time.Second)
		if time.Now().Before(tradeOutEndTime) || time.Now().After(tradeDayOpenEndTime) {
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
				continue
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
	}
}
