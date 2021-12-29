// Package cache package cache
package cache

import (
	"sync"

	"github.com/patrickmn/go-cache"
)

// Cache Cache
type Cache struct {
	Cache *cache.Cache
	lock  sync.RWMutex
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
	newCache.Cache = cache.New(0, 0)
	globalCache = &newCache
}

// Set Set
func (c *Cache) Set(key string, value interface{}) {
	c.lock.Lock()
	c.Cache.Set(key, value, 0)
	c.lock.Unlock()
}
