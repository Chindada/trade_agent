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

// SetTradeDay SetTradeDay
func (c *Cache) SetTradeDay(tradeDay time.Time) {
	key := KeyTradeDay()
	c.getCacheByType(key.Type).Set(key.Name, tradeDay, noExpired)
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

// SetHistroyCloseRange SetHistroyCloseRange
func (c *Cache) SetHistroyCloseRange(closeRange []time.Time) {
	key := KeyHistroyCloseRange()
	c.getCacheByType(key.Type).Set(key.Name, closeRange, noExpired)
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

// SetHistroyTickRange SetHistroyTickRange
func (c *Cache) SetHistroyTickRange(tickRange []time.Time) {
	key := KeyHistroyTickRange()
	c.getCacheByType(key.Type).Set(key.Name, tickRange, noExpired)
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

// SetHistroyKbarRange SetHistroyKbarRange
func (c *Cache) SetHistroyKbarRange(kbarRange []time.Time) {
	key := KeyHistroyKbarRange()
	c.getCacheByType(key.Type).Set(key.Name, kbarRange, noExpired)
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

// KeyIsOpenWithEndWaitTime KeyIsOpenWithEndWaitTime
func KeyIsOpenWithEndWaitTime() *Key {
	return &Key{
		Name: "KeyIsOpenWithEndWaitTime",
		Type: tradeDay,
	}
}

// SetIsOpenWithEndWaitTime SetIsOpenWithEndWaitTime
func (c *Cache) SetIsOpenWithEndWaitTime(isOpen bool) {
	key := KeyIsOpenWithEndWaitTime()
	c.getCacheByType(key.Type).Set(key.Name, isOpen, noExpired)
}

// GetIsOpenWithEndWaitTime GetIsOpenWithEndWaitTime
func (c *Cache) GetIsOpenWithEndWaitTime() bool {
	c.lock.RLock()
	k := KeyIsOpenWithEndWaitTime()
	tmp := c.CacheMap[string(k.Type)]
	c.lock.RUnlock()
	if tmp == nil {
		return false
	}
	if value, ok := tmp.Get(k.Name); ok {
		return value.(bool)
	}
	return false
}

// KeyTradeDayTradeOutEndTime KeyTradeDayTradeOutEndTime
func KeyTradeDayTradeOutEndTime() *Key {
	return &Key{
		Name: "KeyTradeDayTradeOutEndTime",
		Type: tradeDay,
	}
}

// SetTradeDayTradeOutEndTime SetTradeDayTradeOutEndTime
func (c *Cache) SetTradeDayTradeOutEndTime(tradeDayTradeOutEndTime time.Time) {
	key := KeyTradeDayTradeOutEndTime()
	c.getCacheByType(key.Type).Set(key.Name, tradeDayTradeOutEndTime, noExpired)
}

// GetTradeDayTradeOutEndTime GetTradeDayTradeOutEndTime
func (c *Cache) GetTradeDayTradeOutEndTime() time.Time {
	c.lock.RLock()
	k := KeyTradeDayTradeOutEndTime()
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
