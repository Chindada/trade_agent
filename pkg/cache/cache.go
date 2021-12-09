// Package cache package cache
package cache

import (
	"sync"

	"github.com/patrickmn/go-cache"
)

var (
	globalCache *cache.Cache
	once        sync.Once
)

func initGlobalCache() {
	if globalCache != nil {
		return
	}
	globalCache = cache.New(0, 0)
}

// Set Set
func Set(key string, value interface{}) {
	if globalCache == nil {
		once.Do(initGlobalCache)
	}
	globalCache.Set(key, value, 0)
}

// Get Get
func Get(key string) interface{} {
	if globalCache == nil {
		once.Do(initGlobalCache)
	}
	if value, ok := globalCache.Get(key); ok {
		return value
	}
	return nil
}
