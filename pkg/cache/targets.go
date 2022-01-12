// Package cache package cache
package cache

import (
	"trade_agent/pkg/dbagent"
)

// KeyTargets KeyTargets
func KeyTargets() *Key {
	return &Key{
		Name: "KeyTargets",
		Type: targets,
	}
}

// SetTargets SetTargets
func (c *Cache) SetTargets(targetArr []*dbagent.Target) {
	key := KeyTargets()
	c.getCacheByType(key.Type).Set(key.Name, targetArr, noExpired)
}

// GetTargets GetTargets
func (c *Cache) GetTargets() []*dbagent.Target {
	k := KeyTargets()
	if value, ok := c.getCacheByType(k.Type).Get(k.Name); ok {
		return value.([]*dbagent.Target)
	}
	return nil
}

// AppendTargets AppendTargets
func (c *Cache) AppendTargets(arr []*dbagent.Target) {
	tmp := c.GetTargets()
	tmp = append(tmp, arr...)
	c.SetTargets(tmp)
}
