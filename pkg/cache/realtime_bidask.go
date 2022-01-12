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
		Type: realTimeBidaskChannel,
	}
}

// SetRealTimeBidAskChannel SetRealTimeBidAskChannel
func (c *Cache) SetRealTimeBidAskChannel(stockNum string, ch chan *dbagent.RealTimeBidAsk) {
	key := KeyRealTimeBidAskChannel(stockNum)
	c.getCacheByType(key.Type).Set(key.Name, ch, noExpired)
}

// GetRealTimeBidAskChannel GetRealTimeBidAskChannel
func (c *Cache) GetRealTimeBidAskChannel(stockNum string) chan *dbagent.RealTimeBidAsk {
	k := KeyRealTimeBidAskChannel(stockNum)
	if value, ok := c.getCacheByType(k.Type).Get(k.Name); ok {
		return value.(chan *dbagent.RealTimeBidAsk)
	}
	return nil
}

// KeyRealTimeBidAskStatus KeyRealTimeBidAskStatus
func KeyRealTimeBidAskStatus(stockNum string) *Key {
	return &Key{
		Name: fmt.Sprintf("KeyRealTimeBidAskStatus:%s", stockNum),
		Type: keyTypeRealTimeBidAskStatus(stockNum),
	}
}

// SetRealTimeBidAskStatus SetRealTimeBidAskStatus
func (c *Cache) SetRealTimeBidAskStatus(stockNum string, status string) {
	key := KeyRealTimeBidAskStatus(stockNum)
	c.getCacheByType(key.Type).Set(key.Name, status, noExpired)
}

// GetRealTimeBidAskStatus GetRealTimeBidAskStatus
func (c *Cache) GetRealTimeBidAskStatus(stockNum string) string {
	k := KeyRealTimeBidAskStatus(stockNum)
	if value, ok := c.getCacheByType(k.Type).Get(k.Name); ok {
		return value.(string)
	}
	return ""
}
