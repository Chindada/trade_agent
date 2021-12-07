// Package cache package cache
package cache

import (
	"sync"

	"github.com/patrickmn/go-cache"
)

var (
	globalCache *cache.Cache
	initLock    sync.Mutex
)

func initGlobalCache() {
	defer initLock.Unlock()
	initLock.Lock()
	if globalCache != nil {
		return
	}
	globalCache = cache.New(0, 0)
}

// Set Set
func Set(key string, value interface{}) {
	if globalCache == nil {
		initGlobalCache()
	}
	globalCache.Set(key, value, 0)
}

// Get Get
func Get(key string) interface{} {
	if globalCache == nil {
		initGlobalCache()
	}
	if value, ok := globalCache.Get(key); ok {
		return value
	}
	return nil
}
