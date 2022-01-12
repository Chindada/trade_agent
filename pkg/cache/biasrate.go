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

// SetBiasRate SetBiasRate
func (c *Cache) SetBiasRate(stockNum string, biasRate float64) {
	key := KeyBiasRate(stockNum)
	c.getCacheByType(key.Type).Set(key.Name, biasRate, noExpired)
}

// GetBiasRate GetBiasRate
func (c *Cache) GetBiasRate(stockNum string) float64 {
	k := KeyBiasRate(stockNum)
	if value, ok := c.getCacheByType(k.Type).Get(k.Name); ok {
		return value.(float64)
	}
	return 0
}
