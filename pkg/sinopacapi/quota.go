// Package sinopacapi package sinopacapi
package sinopacapi

import (
	"math"
)

var (
	tradeTaxRatio float64
	tradeFeeRatio float64
	feeDiscount   float64
)

// SetOrderToQuota SetOrderToQuota
func (c *TradeAgent) SetOrderToQuota(order Order, success bool) {
	defer c.mu.Unlock()
	c.mu.Lock()
	var cost int64
	switch order.Action {
	case ActionBuy:
		cost = GetStockBuyCost(order.Price, order.Quantity)
	case ActionSellFirst:
		cost = GetStockSellCost(order.Price, order.Quantity)
	}

	if success {
		cost = -cost
	}
	c.tradeQuota += cost
}

// GetStockBuyCost GetStockBuyCost
func GetStockBuyCost(price float64, qty int64) int64 {
	return int64(math.Ceil(price*float64(qty)*1000) + math.Floor(price*float64(qty)*1000*tradeFeeRatio))
}

// GetStockSellCost GetStockSellCost
func GetStockSellCost(price float64, qty int64) int64 {
	return int64(math.Ceil(price*float64(qty)*1000) - math.Floor(price*float64(qty)*1000*tradeFeeRatio) - math.Floor(price*float64(qty)*1000*tradeTaxRatio))
}

// GetStockTradeFeeDiscount GetStockTradeFeeDiscount
func GetStockTradeFeeDiscount(price float64, qty int64) int64 {
	return int64(math.Floor(price*float64(qty)*1000*tradeFeeRatio) * (1 - feeDiscount))
}
