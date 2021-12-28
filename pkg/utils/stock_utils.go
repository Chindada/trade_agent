// Package utils package utils
package utils

import "math"

// TradeQuota TradeQuota
var TradeQuota int64 = 1000000

const (
	// TradeFeeRatio TradeFeeRatio
	TradeFeeRatio float64 = 0.001425
	// FeeDiscount FeeDiscount
	FeeDiscount float64 = 0.35
	// TradeTaxRatio TradeTaxRatio
	TradeTaxRatio float64 = 0.0015
)

// GetStockBuyCost GetStockBuyCost
func GetStockBuyCost(price float64, qty int64) int64 {
	return int64(math.Ceil(price*float64(qty)*1000) + math.Floor(price*float64(qty)*1000*TradeFeeRatio))
}

// GetStockSellCost GetStockSellCost
func GetStockSellCost(price float64, qty int64) int64 {
	return int64(math.Ceil(price*float64(qty)*1000) - math.Floor(price*float64(qty)*1000*TradeFeeRatio) - math.Floor(price*float64(qty)*1000*TradeTaxRatio))
}

// GetStockTradeFeeDiscount GetStockTradeFeeDiscount
func GetStockTradeFeeDiscount(price float64, qty int64) int64 {
	return int64(math.Floor(price*float64(qty)*1000*TradeFeeRatio) * (1 - FeeDiscount))
}

// GetNewClose GetNewClose
func GetNewClose(close float64, unit int64) float64 {
	if close == 0 {
		return 0
	}
	for {
		if unit == 0 {
			return Round(close, 2)
		}
		diff := GetDiff(close)
		if unit > 0 {
			close += diff
			unit--
		} else {
			close -= diff
			unit++
		}
	}
}

// GetMaxByOpen GetMaxByOpen
func GetMaxByOpen(open float64) float64 {
	if open == 0 {
		return 0
	}
	tmpClose := open
	var changeRate, diff float64
	for {
		diff = GetDiff(tmpClose)
		tmpClose += diff
		tmpClose = Round(tmpClose, 2)
		changeRate = 100 * (tmpClose - open) / open
		if Round(changeRate, 2) > 10 {
			return Round(tmpClose-diff, 2)
		}
	}
}

// GetMinByOpen GetMinByOpen
func GetMinByOpen(open float64) float64 {
	if open == 0 {
		return 0
	}
	tmpClose := open
	var changeRate, diff float64
	for {
		diff = GetDiff(tmpClose)
		tmpClose -= diff
		tmpClose = Round(tmpClose, 2)
		changeRate = 100 * (tmpClose - open) / open
		if Round(changeRate, 2) < -10 {
			return Round(tmpClose+diff, 2)
		}
	}
}

// GetDiff GetDiff
func GetDiff(close float64) float64 {
	switch {
	case close > 0 && close < 10:
		return 0.01
	case close >= 10 && close < 50:
		return 0.05
	case close >= 50 && close < 100:
		return 0.1
	case close >= 100 && close < 500:
		return 0.5
	case close >= 500 && close < 1000:
		return 1
	case close >= 1000:
		return 5
	}
	return 0
}
