// Package cache package cache
package cache

import (
	"fmt"
	"time"
	"trade_agent/global"
)

// KeyStockHistoryClose KeyStockHistoryClose
func KeyStockHistoryClose(stockNum string, date time.Time) string {
	return fmt.Sprintf("KeyStockHistoryClose:%s:%s", stockNum, date.Format(global.ShortTimeLayout))
}

// GetHistoryClose GetHistoryClose
func (c *Cache) GetHistoryClose(stockNum string, date time.Time) float64 {
	defer c.lock.RUnlock()
	c.lock.RLock()
	if value, ok := c.Cache.Get(KeyStockHistoryClose(stockNum, date)); ok {
		return value.(float64)
	}
	return 0
}
