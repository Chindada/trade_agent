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
