// Package pb package pb
package pb

import "trade_agent/pkg/dbagent"

// ToStock ToStock
func (c *StockMessage) ToStock() *dbagent.Stock {
	var dayTrade bool
	if c.GetDayTrade() == "Yes" {
		dayTrade = true
	}
	return &dbagent.Stock{
		Number:             c.GetCode(),
		Name:               c.GetName(),
		Exchange:           c.GetExchange(),
		Category:           c.GetCategory(),
		DayTrade:           dayTrade,
		LastClose:          c.GetReference(),
		LastVolume:         0,
		LastCloseChangePct: 0,
	}
}
