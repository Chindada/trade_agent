// Package cache package cache
package cache

import (
	"fmt"
	"trade_agent/pkg/dbagent"
)

// KeyRealTimeBidAskChannel KeyRealTimeBidAskChannel
func KeyRealTimeBidAskChannel(stockNum string) *Key {
	return &Key{
		Name: fmt.Sprintf("KeyRealTimeBidAskChannel:%s", stockNum),
		Type: realTimeBidask,
	}
}

// SetRealTimeBidAskChannel SetRealTimeBidAskChannel
func (c *Cache) SetRealTimeBidAskChannel(stockNum string, ch chan *dbagent.RealTimeBidAsk) {
	key := KeyRealTimeBidAskChannel(stockNum)
	c.getCacheByType(key.Type).Set(key.Name, ch, noExpired)
}

// GetRealTimeBidAskChannel GetRealTimeBidAskChannel
func (c *Cache) GetRealTimeBidAskChannel(stockNum string) chan *dbagent.RealTimeBidAsk {
	c.lock.RLock()
	k := KeyRealTimeBidAskChannel(stockNum)
	tmp := c.CacheMap[string(k.Type)]
	c.lock.RUnlock()
	if tmp == nil {
		return nil
	}
	if value, ok := tmp.Get(k.Name); ok {
		return value.(chan *dbagent.RealTimeBidAsk)
	}
	return nil
}

// KeyRealTimeBidAskStatus KeyRealTimeBidAskStatus
func KeyRealTimeBidAskStatus(stockNum string) *Key {
	return &Key{
		Name: fmt.Sprintf("KeyRealTimeBidAskStatus:%s", stockNum),
		Type: realTimeBidask,
	}
}

// SetRealTimeBidAskStatus SetRealTimeBidAskStatus
func (c *Cache) SetRealTimeBidAskStatus(stockNum string, status string) {
	key := KeyRealTimeBidAskStatus(stockNum)
	c.getCacheByType(key.Type).Set(key.Name, status, noExpired)
}

// GetRealTimeBidAskStatus GetRealTimeBidAskStatus
func (c *Cache) GetRealTimeBidAskStatus(stockNum string) string {
	c.lock.RLock()
	k := KeyRealTimeBidAskStatus(stockNum)
	tmp := c.CacheMap[string(k.Type)]
	c.lock.RUnlock()
	if tmp == nil {
		return ""
	}
	if value, ok := tmp.Get(k.Name); ok {
		return value.(string)
	}
	return ""
}
