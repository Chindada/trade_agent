// Package cache package cache
package cache

import (
	"fmt"
	"time"
	"trade_agent/global"
)

// KeyStockHistoryClose KeyStockHistoryClose
func KeyStockHistoryClose(stockNum string, date time.Time) string {
	return fmt.Sprintf("HistoryClose:%s:%s", stockNum, date.Format(global.ShortTimeLayout))
}
