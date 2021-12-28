// Package cache package cache
package cache

import (
	"fmt"
)

// KeyStockHistoryKbarAnalyze KeyStockHistoryKbarAnalyze
func KeyStockHistoryKbarAnalyze(stockNum string) string {
	return fmt.Sprintf("KeyStockHistoryKbarAnalyze:%s", stockNum)
}

// GetStockHistoryKbarAnalyze GetStockHistoryKbarAnalyze
func (c *Cache) GetStockHistoryKbarAnalyze(stockNum string) string {
	defer c.lock.RUnlock()
	c.lock.RLock()
	if value, ok := c.Cache.Get(KeyStockHistoryKbarAnalyze(stockNum)); ok {
		return value.(string)
	}
	return ""
}
