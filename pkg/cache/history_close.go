// Package cache package cache
package cache

import (
	"fmt"
	"time"
	"trade_agent/global"
)

// KeyStockHistoryClose KeyStockHistoryClose
func KeyStockHistoryClose(stockNum string, date time.Time) *Key {
	return &Key{
		Name: fmt.Sprintf("KeyStockHistoryClose:%s:%s", stockNum, date.Format(global.ShortTimeLayout)),
		Type: historyClose,
	}
}

// SetStockHistoryClose SetStockHistoryClose
func (c *Cache) SetStockHistoryClose(stockNum string, close float64, date time.Time) {
	key := KeyStockHistoryClose(stockNum, date)
	c.getCacheByType(key.Type).Set(key.Name, close, noExpired)
}

// GetHistoryClose GetHistoryClose
func (c *Cache) GetHistoryClose(stockNum string, date time.Time) float64 {
	k := KeyStockHistoryClose(stockNum, date)
	if value, ok := c.getCacheByType(k.Type).Get(k.Name); ok {
		return value.(float64)
	}
	return 0
}
