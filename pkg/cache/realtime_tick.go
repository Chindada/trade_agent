// Package cache package cache
package cache

import (
	"fmt"
	"trade_agent/pkg/dbagent"
)

// KeyRealTimeTickChannel KeyRealTimeTickChannel
func KeyRealTimeTickChannel(stockNum string) *Key {
	return &Key{
		Name: fmt.Sprintf("KeyRealTimeTickChannel:%s", stockNum),
		Type: realTimeTickChannel,
	}
}

// SetRealTimeTickChannel SetRealTimeTickChannel
func (c *Cache) SetRealTimeTickChannel(stockNum string, ch chan *dbagent.RealTimeTick) {
	key := KeyRealTimeTickChannel(stockNum)
	c.getCacheByType(key.Type).Set(key.Name, ch, noExpired)
}

// GetRealTimeTickChannel GetRealTimeTickChannel
func (c *Cache) GetRealTimeTickChannel(stockNum string) chan *dbagent.RealTimeTick {
	c.lock.RLock()
	k := KeyRealTimeTickChannel(stockNum)
	tmp := c.CacheMap[string(k.Type)]
	c.lock.RUnlock()
	if tmp == nil {
		return nil
	}
	if value, ok := tmp.Get(k.Name); ok {
		return value.(chan *dbagent.RealTimeTick)
	}
	return nil
}

// KeyRealTimeTickClose KeyRealTimeTickClose
func KeyRealTimeTickClose(stockNum string) *Key {
	return &Key{
		Name: fmt.Sprintf("KeyRealTimeTickClose:%s", stockNum),
		Type: keyTypeRealTimeClose(stockNum),
	}
}

// SetRealTimeTickClose SetRealTimeTickClose
func (c *Cache) SetRealTimeTickClose(stockNum string, close float64) {
	key := KeyRealTimeTickClose(stockNum)
	c.getCacheByType(key.Type).Set(key.Name, close, noExpired)
}

// GetRealTimeTickClose GetRealTimeTickClose
func (c *Cache) GetRealTimeTickClose(stockNum string) float64 {
	c.lock.RLock()
	k := KeyRealTimeTickClose(stockNum)
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
