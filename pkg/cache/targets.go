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

// GetTargets GetTargets
func (c *Cache) GetTargets() []*dbagent.Target {
	c.lock.RLock()
	k := KeyTargets()
	tmp := c.CacheMap[string(k.Type)]
	c.lock.RUnlock()
	if tmp == nil {
		return nil
	}
	if value, ok := tmp.Get(k.Name); ok {
		return value.([]*dbagent.Target)
	}
	return nil
}

// AppendTargets AppendTargets
func (c *Cache) AppendTargets(arr []*dbagent.Target) {
	tmp := c.GetTargets()
	tmp = append(tmp, arr...)
	c.Set(KeyTargets(), tmp)
}
