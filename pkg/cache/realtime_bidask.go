// Package cache package cache
package cache

import (
	"fmt"
	"trade_agent/pkg/dbagent"
)

// KeyRealTimeBidAskChannel KeyRealTimeBidAskChannel
func KeyRealTimeBidAskChannel(stockNum string) string {
	return fmt.Sprintf("KeyRealTimeBidAskChannel:%s", stockNum)
}

// KeyRealTimeBidAskStatus KeyRealTimeBidAskStatus
func KeyRealTimeBidAskStatus(stockNum string) string {
	return fmt.Sprintf("KeyRealTimeBidAskStatus:%s", stockNum)
}

// GetRealTimeBidAskChannel GetRealTimeBidAskChannel
func (c *Cache) GetRealTimeBidAskChannel(stockNum string) chan *dbagent.RealTimeBidAsk {
	defer c.lock.RUnlock()
	c.lock.RLock()
	if value, ok := c.Cache.Get(KeyRealTimeBidAskChannel(stockNum)); ok {
		return value.(chan *dbagent.RealTimeBidAsk)
	}
	return nil
}

// GetRealTimeBidAskStatus GetRealTimeBidAskStatus
func (c *Cache) GetRealTimeBidAskStatus(stockNum string) string {
	defer c.lock.RUnlock()
	c.lock.RLock()
	if value, ok := c.Cache.Get(KeyRealTimeBidAskChannel(stockNum)); ok {
		return value.(string)
	}
	return ""
}
