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
	k := KeyStockHistoryKbarAnalyze(stockNum)
	if value, ok := c.getCacheByType(k.Type).Get(k.Name); ok {
		return value.(string)
	}
	return ""
}
