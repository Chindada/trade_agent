// Package cache package cache
package cache

import (
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
)

const (
	noExpired time.Duration = 0
	noCleanUp time.Duration = 0
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

// getCacheByType getCacheByType
func (c *Cache) getCacheByType(keyType keyType) *cache.Cache {
	c.lock.RLock()
	tmp := c.CacheMap[string(keyType)]
	c.lock.RUnlock()
	if tmp == nil {
		tmp = cache.New(noExpired, noCleanUp)
		c.lock.Lock()
		c.CacheMap[string(keyType)] = tmp
		c.lock.Unlock()
	}
	return tmp
}
