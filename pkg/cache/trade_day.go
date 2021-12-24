// Package cache package cache
package cache

import (
	"time"
)

// KeyTradeDay KeyTradeDay
func KeyTradeDay() string {
	return "TradeDay"
}

// KeyHistroyCloseRange KeyHistroyCloseRange
func KeyHistroyCloseRange() string {
	return "KeyHistroyCloseRange"
}

// KeyHistroyTickRange KeyHistroyTickRange
func KeyHistroyTickRange() string {
	return "KeyHistroyTickRange"
}

// KeyHistroyKbarRange KeyHistroyKbarRange
func KeyHistroyKbarRange() string {
	return "KeyHistroyKbarRange"
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

// GetHistroyCloseRange GetHistroyCloseRange
func (c *Cache) GetHistroyCloseRange() []time.Time {
	defer c.lock.RUnlock()
	c.lock.RLock()
	if value, ok := c.Cache.Get(KeyHistroyCloseRange()); ok {
		return value.([]time.Time)
	}
	return []time.Time{}
}

// GetHistroyTickRange GetHistroyTickRange
func (c *Cache) GetHistroyTickRange() []time.Time {
	defer c.lock.RUnlock()
	c.lock.RLock()
	if value, ok := c.Cache.Get(KeyHistroyTickRange()); ok {
		return value.([]time.Time)
	}
	return []time.Time{}
}

// GetHistroyKbarRange GetHistroyKbarRange
func (c *Cache) GetHistroyKbarRange() []time.Time {
	defer c.lock.RUnlock()
	c.lock.RLock()
	if value, ok := c.Cache.Get(KeyHistroyKbarRange()); ok {
		return value.([]time.Time)
	}
	return []time.Time{}
}
