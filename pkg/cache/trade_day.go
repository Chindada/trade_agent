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
	k := KeyTradeDay()
	if value, ok := c.getCacheByType(k.Type).Get(k.Name); ok {
		return value.(time.Time)
	}
	return time.Time{}
}

// KeyTradeDayOpenEndTime KeyTradeDayOpenEndTime
func KeyTradeDayOpenEndTime() *Key {
	return &Key{
		Name: "KeyTradeDayOpenEndTime",
		Type: tradeDay,
	}
}

// SetTradeDayOpenEndTime SetTradeDayOpenEndTime
func (c *Cache) SetTradeDayOpenEndTime(tradeDay time.Time) {
	key := KeyTradeDayOpenEndTime()
	c.getCacheByType(key.Type).Set(key.Name, tradeDay, noExpired)
}

// GetTradeDayOpenEndTime GetTradeDayOpenEndTime
func (c *Cache) GetTradeDayOpenEndTime() time.Time {
	k := KeyTradeDayOpenEndTime()
	if value, ok := c.getCacheByType(k.Type).Get(k.Name); ok {
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
	k := KeyHistroyCloseRange()
	if value, ok := c.getCacheByType(k.Type).Get(k.Name); ok {
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
	k := KeyHistroyTickRange()
	if value, ok := c.getCacheByType(k.Type).Get(k.Name); ok {
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
	k := KeyHistroyKbarRange()
	if value, ok := c.getCacheByType(k.Type).Get(k.Name); ok {
		return value.([]time.Time)
	}
	return []time.Time{}
}

// KeyIsAllowTrade KeyIsAllowTrade
func KeyIsAllowTrade() *Key {
	return &Key{
		Name: "KeyIsAllowTrade",
		Type: tradeDay,
	}
}

// SetIsAllowTrade SetIsAllowTrade
func (c *Cache) SetIsAllowTrade(isAllowTrade bool) {
	key := KeyIsAllowTrade()
	c.getCacheByType(key.Type).Set(key.Name, isAllowTrade, noExpired)
}

// GetIsAllowTrade GetIsAllowTrade
func (c *Cache) GetIsAllowTrade() bool {
	k := KeyIsAllowTrade()
	if value, ok := c.getCacheByType(k.Type).Get(k.Name); ok {
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
	k := KeyTradeDayTradeOutEndTime()
	if value, ok := c.getCacheByType(k.Type).Get(k.Name); ok {
		return value.(time.Time)
	}
	return time.Time{}
}
