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
		Type: realTimeTick,
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
