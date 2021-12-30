// Package cache package cache
package cache

import (
	"fmt"
)

// KeyBiasRate KeyBiasRate
func KeyBiasRate(stockNum string) *Key {
	return &Key{
		Name: fmt.Sprintf("KeyBiasRate:%s", stockNum),
		Type: biasRate,
	}
}

// GetBiasRate GetBiasRate
func (c *Cache) GetBiasRate(stockNum string) float64 {
	c.lock.RLock()
	k := KeyBiasRate(stockNum)
	tmp := c.CacheMap[string(k.Type)]
	c.lock.RUnlock()
	if tmp == nil {
		return 0
	}
	if value, ok := tmp.Get(k.Name); ok {
		return value.(float64)
	}
	return 0
}
