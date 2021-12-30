// Package cache package cache
package cache

import (
	"time"
)

// KeyTradeDay KeyTradeDay
func KeyTradeDay() *Key {
	return &Key{
		Name: "KeyTradeDay",
		Type: tradeDay,
	}
}

// GetTradeDay GetTradeDay
func (c *Cache) GetTradeDay() time.Time {
	c.lock.RLock()
	k := KeyTradeDay()
	tmp := c.CacheMap[string(k.Type)]
	c.lock.RUnlock()
	if tmp == nil {
		return time.Time{}
	}
	if value, ok := tmp.Get(k.Name); ok {
		return value.(time.Time)
	}
	return time.Time{}
}

// KeyHistroyCloseRange KeyHistroyCloseRange
func KeyHistroyCloseRange() *Key {
	return &Key{
		Name: "KeyHistroyCloseRange",
		Type: tradeDay,
	}
}

// GetHistroyCloseRange GetHistroyCloseRange
func (c *Cache) GetHistroyCloseRange() []time.Time {
	c.lock.RLock()
	k := KeyHistroyCloseRange()
	tmp := c.CacheMap[string(k.Type)]
	c.lock.RUnlock()
	if tmp == nil {
		return []time.Time{}
	}
	if value, ok := tmp.Get(k.Name); ok {
		return value.([]time.Time)
	}
	return []time.Time{}
}

// KeyHistroyTickRange KeyHistroyTickRange
func KeyHistroyTickRange() *Key {
	return &Key{
		Name: "KeyHistroyTickRange",
		Type: tradeDay,
	}
}

// GetHistroyTickRange GetHistroyTickRange
func (c *Cache) GetHistroyTickRange() []time.Time {
	c.lock.RLock()
	k := KeyHistroyTickRange()
	tmp := c.CacheMap[string(k.Type)]
	c.lock.RUnlock()
	if tmp == nil {
		return []time.Time{}
	}
	if value, ok := tmp.Get(k.Name); ok {
		return value.([]time.Time)
	}
	return []time.Time{}
}

// KeyHistroyKbarRange KeyHistroyKbarRange
func KeyHistroyKbarRange() *Key {
	return &Key{
		Name: "KeyHistroyKbarRange",
		Type: tradeDay,
	}
}

// GetHistroyKbarRange GetHistroyKbarRange
func (c *Cache) GetHistroyKbarRange() []time.Time {
	c.lock.RLock()
	k := KeyHistroyKbarRange()
	tmp := c.CacheMap[string(k.Type)]
	c.lock.RUnlock()
	if tmp == nil {
		return []time.Time{}
	}
	if value, ok := tmp.Get(k.Name); ok {
		return value.([]time.Time)
	}
	return []time.Time{}
}
