// Package cache package cache
package cache

import "trade_agent/pkg/dbagent"

// KeyTSESnapshot KeyTSESnapshot
func KeyTSESnapshot() *Key {
	return &Key{
		Name: "KeyTSESnapshot",
		Type: tseSnapshot,
	}
}

// SetTSESnapshot SetTSESnapshot
func (c *Cache) SetTSESnapshot(snapshot *dbagent.TSESnapShot) {
	key := KeyTargets()
	c.getCacheByType(key.Type).Set(key.Name, snapshot, noExpired)
}

// GetTSESnapshot GetTSESnapshot
func (c *Cache) GetTSESnapshot() *dbagent.TSESnapShot {
	k := KeyTargets()
	if value, ok := c.getCacheByType(k.Type).Get(k.Name); ok {
		return value.(*dbagent.TSESnapShot)
	}
	return nil
}
