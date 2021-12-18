// Package cache package cache
package cache

import (
	"time"
)

// KeyTradeDay KeyTradeDay
func KeyTradeDay() string {
	return "TradeDay"
}

// KeyHistroyRange KeyHistroyRange
func KeyHistroyRange() string {
	return "KeyHistroyRange"
}

// GetTradeDay GetTradeDay
func (c *Cache) GetTradeDay() time.Time {
	defer c.lock.RUnlock()
	c.lock.RLock()
	if value, ok := c.Cache.Get(KeyTradeDay()); ok {
		return value.(time.Time)
	}
	return time.Time{}
}

// GetHistroyRange GetHistroyRange
func (c *Cache) GetHistroyRange() []time.Time {
	defer c.lock.RUnlock()
	c.lock.RLock()
	if value, ok := c.Cache.Get(KeyHistroyRange()); ok {
		return value.([]time.Time)
	}
	return []time.Time{}
}
