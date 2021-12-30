// Package cache package cache
package cache

import (
	"fmt"
)

// KeyStockHistoryTickAnalyze KeyStockHistoryTickAnalyze
func KeyStockHistoryTickAnalyze(stockNum string) *Key {
	return &Key{
		Name: fmt.Sprintf("KeyStockHistoryTickAnalyze:%s", stockNum),
		Type: historyTick,
	}
}

// GetStockHistoryTickAnalyze GetStockHistoryTickAnalyze
func (c *Cache) GetStockHistoryTickAnalyze(stockNum string) []int64 {
	c.lock.RLock()
	k := KeyStockHistoryTickAnalyze(stockNum)
	tmp := c.CacheMap[string(k.Type)]
	c.lock.RUnlock()
	if tmp == nil {
		return []int64{}
	}
	if value, ok := tmp.Get(k.Name); ok {
		return value.([]int64)
	}
	return []int64{}
}
