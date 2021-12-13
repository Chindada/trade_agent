// Package cache package cache
package cache

import (
	"fmt"
	"time"
	"trade_agent/global"
)

// KeyStockDetail KeyStockDetail
func KeyStockDetail(stockNum string) string {
	return fmt.Sprintf("StockDetail:%s", stockNum)
}

// KeyStockHistoryClose KeyStockHistoryClose
func KeyStockHistoryClose(stockNum string, date time.Time) string {
	return fmt.Sprintf("HistoryClose:%s:%s", stockNum, date.Format(global.ShortTimeLayout))
}
