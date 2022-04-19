// Package cache package cache
package cache

import (
	"fmt"

	"trade_agent/pkg/dbagent"
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

// KeyStockHistoryDayKbar KeyStockHistoryDayKbar
func KeyStockHistoryDayKbar(date, stockNum string) *Key {
	return &Key{
		Name: fmt.Sprintf("KeyStockHistoryDayKbar:%s:%s", date, stockNum),
		Type: keyTypeHistoryKbarAnalyze(stockNum),
	}
}

// SetStockHistoryDayKbar SetStockHistoryDayKbar
func (c *Cache) SetStockHistoryDayKbar(date, stockNum string, kbar *dbagent.HistoryKbar) {
	key := KeyStockHistoryDayKbar(date, stockNum)
	c.getCacheByType(key.Type).Set(key.Name, kbar, noExpired)
}

// GetStockHistoryDayKbar GetStockHistoryDayKbar
func (c *Cache) GetStockHistoryDayKbar(date, stockNum string) *dbagent.HistoryKbar {
	k := KeyStockHistoryDayKbar(date, stockNum)
	if value, ok := c.getCacheByType(k.Type).Get(k.Name); ok {
		return value.(*dbagent.HistoryKbar)
	}
	return nil
}
