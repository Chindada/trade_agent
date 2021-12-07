// Package cache package cache
package cache

import (
	"time"

	"gitlab.tocraw.com/root/toc_trader/global"
)

// KeyStockHistoryClose KeyStockHistoryClose
func KeyStockHistoryClose(stockNum string, date time.Time) string {
	return "HistoryClose" + ":" + stockNum + ":" + date.Format(global.ShortTimeLayout)
}
