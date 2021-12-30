// Package cache package cache
package cache

import (
	"sync"

	"github.com/patrickmn/go-cache"
)

// Cache Cache
type Cache struct {
	CacheMap map[string]*cache.Cache
	lock     sync.RWMutex
}

// Key Key
type Key struct {
	Name string  `json:"name,omitempty" yaml:"name"`
	Type keyType `json:"type,omitempty" yaml:"type"`
}

var (
	globalCache *Cache
	once        sync.Once
)

// GetCache GetCache
func GetCache() *Cache {
	if globalCache != nil {
		return globalCache
	}
	once.Do(initGlobalCache)
	return globalCache
}

func initGlobalCache() {
	if globalCache != nil {
		return
	}
	var newCache Cache
	newCache.CacheMap = make(map[string]*cache.Cache)
	globalCache = &newCache
}

// Set Set
func (c *Cache) Set(key *Key, value interface{}) {
	c.lock.RLock()
	tmp := c.CacheMap[string(key.Type)]
	c.lock.RUnlock()
	if tmp == nil {
		tmp = cache.New(0, 0)
		c.lock.Lock()
		c.CacheMap[string(key.Type)] = tmp
		c.lock.Unlock()
	}
	tmp.Set(key.Name, value, 0)
}
