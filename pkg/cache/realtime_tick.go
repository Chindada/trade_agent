// Package cache package cache
package cache

import (
	"fmt"
	"trade_agent/pkg/dbagent"
)

// KeyRealTimeTickChannel KeyRealTimeTickChannel
func KeyRealTimeTickChannel(stockNum string) string {
	return fmt.Sprintf("KeyRealTimeTickChannel:%s", stockNum)
}

// GetRealTimeTickChannel GetRealTimeTickChannel
func (c *Cache) GetRealTimeTickChannel(stockNum string) chan *dbagent.RealTimeTick {
	defer c.lock.RUnlock()
	c.lock.RLock()
	if value, ok := c.Cache.Get(KeyRealTimeTickChannel(stockNum)); ok {
		return value.(chan *dbagent.RealTimeTick)
	}
	return nil
}
