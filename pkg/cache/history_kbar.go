// Package cache package cache
package cache

import (
	"fmt"
)

// KeyStockHistoryKbarAnalyze KeyStockHistoryKbarAnalyze
func KeyStockHistoryKbarAnalyze(stockNum string) *Key {
	return &Key{
		Name: fmt.Sprintf("KeyStockHistoryKbarAnalyze:%s", stockNum),
		Type: keyTypeHistoryKbarAnalyze(stockNum),
	}
}

// SetStockHistoryKbarAnalyze SetStockHistoryKbarAnalyze
func (c *Cache) SetStockHistoryKbarAnalyze(stockNum string, status string) {
	key := KeyStockHistoryKbarAnalyze(stockNum)
	c.getCacheByType(key.Type).Set(key.Name, status, noExpired)
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
