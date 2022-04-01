// Package cache package cache
package cache

import (
	"fmt"
)

// KeyStockVolume KeyStockVolume
func KeyStockVolume(stockNum string) *Key {
	return &Key{
		Name: fmt.Sprintf("KeyStockVolume:%s", stockNum),
		Type: historyClose,
	}
}

// SetStockVolume SetStockVolume
func (c *Cache) SetStockVolume(stockNum string, volume int64) {
	key := KeyStockVolume(stockNum)
	c.getCacheByType(key.Type).Set(key.Name, volume, noExpired)
}

// GetStockVolume GetStockVolume
func (c *Cache) GetStockVolume(stockNum string) int64 {
	k := KeyStockVolume(stockNum)
	if value, ok := c.getCacheByType(k.Type).Get(k.Name); ok {
		return value.(int64)
	}
	return 0
}
