// Package cache package cache
package cache

import (
	"fmt"
)

// KeyBiasRate KeyBiasRate
func KeyBiasRate(stockNum string) string {
	return fmt.Sprintf("KeyBiasRate:%s", stockNum)
}

// GetBiasRate GetBiasRate
func (c *Cache) GetBiasRate(stockNum string) float64 {
	defer c.lock.RUnlock()
	c.lock.RLock()
	if value, ok := c.Cache.Get(KeyBiasRate(stockNum)); ok {
		return value.(float64)
	}
	return 0
}
