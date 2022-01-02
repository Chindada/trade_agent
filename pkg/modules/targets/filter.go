// Package targets package targets
package targets

import (
	"trade_agent/pkg/config"
	"trade_agent/pkg/dbagent"
)

type stockWithData struct {
	stock       *dbagent.Stock
	totalVolume int64
	close       float64
}

func stockTargetFilter(v stockWithData, cond config.TargetCond, isRealTime bool) bool {
	blackCategoryMap := make(map[string]bool)
	blackStockMap := make(map[string]bool)
	for _, v := range cond.BlackCategory {
		blackCategoryMap[v] = true
	}
	for _, v := range cond.BlackStock {
		blackStockMap[v] = true
	}

	if blackStockMap[v.stock.Number] {
		return false
	}
	if blackCategoryMap[v.stock.Category] {
		return false
	}
	if v.totalVolume < cond.LimitVolume && !isRealTime {
		return false
	}
	if v.close < cond.LimitPriceLow || v.close > cond.LimitPriceHigh {
		return false
	}
	return true
}
