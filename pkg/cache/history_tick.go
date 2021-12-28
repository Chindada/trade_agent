// Package cache package cache
package cache

import (
	"fmt"
)

// KeyStockHistoryTickAnalyze KeyStockHistoryTickAnalyze
func KeyStockHistoryTickAnalyze(stockNum string) string {
	return fmt.Sprintf("KeyStockHistoryTickAnalyze:%s", stockNum)
}

// GetStockHistoryTickAnalyze GetStockHistoryTickAnalyze
func (c *Cache) GetStockHistoryTickAnalyze(stockNum string) []int64 {
	defer c.lock.RUnlock()
	c.lock.RLock()
	if value, ok := c.Cache.Get(KeyStockHistoryTickAnalyze(stockNum)); ok {
		return value.([]int64)
	}
	return []int64{}
}
