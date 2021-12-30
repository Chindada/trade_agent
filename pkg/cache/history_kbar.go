// Package cache package cache
package cache

import (
	"fmt"
)

// KeyStockHistoryKbarAnalyze KeyStockHistoryKbarAnalyze
func KeyStockHistoryKbarAnalyze(stockNum string) *Key {
	return &Key{
		Name: fmt.Sprintf("KeyStockHistoryKbarAnalyze:%s", stockNum),
		Type: historyKbar,
	}
}

// GetStockHistoryKbarAnalyze GetStockHistoryKbarAnalyze
func (c *Cache) GetStockHistoryKbarAnalyze(stockNum string) string {
	c.lock.RLock()
	k := KeyStockHistoryKbarAnalyze(stockNum)
	tmp := c.CacheMap[string(k.Type)]
	c.lock.RUnlock()
	if tmp == nil {
		return ""
	}
	if value, ok := tmp.Get(k.Name); ok {
		return value.(string)
	}
	return ""
}
