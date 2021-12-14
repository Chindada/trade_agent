// Package cache package cache
package cache

import (
	"sync"

	"github.com/patrickmn/go-cache"
)

// Cache Cache
type Cache struct {
	Cache *cache.Cache
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
	c.Cache.Set(key, value, 0)
}

// Get Get
func (c *Cache) Get(key string) interface{} {
	if value, ok := c.Cache.Get(key); ok {
		return value
	}
	return nil
}
