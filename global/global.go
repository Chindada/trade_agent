// Package global package global
package global

import "sync"

const (
	// LongTimeLayout LongTimeLayout
	LongTimeLayout string = "2006-01-02 15:04:05"
	// ShortTimeLayout ShortTimeLayout
	ShortTimeLayout string = "2006-01-02"
)

// Setting Setting
type Setting struct {
	lock          sync.RWMutex
	basePath      string
	isDevelopment bool
}

var globalSetting *Setting

// Get Get
func Get() *Setting {
	if globalSetting != nil {
		return globalSetting
	}
	globalSetting = &Setting{}
	return globalSetting
}

// SetBasePath SetBasePath
func (c *Setting) SetBasePath(path string) {
	defer c.lock.RUnlock()
	c.lock.RLock()
	c.basePath = path
}

// GetBasePath GetBasePath
func (c *Setting) GetBasePath() string {
	defer c.lock.RUnlock()
	c.lock.RLock()
	return c.basePath
}

// SetIsDevelopment SetIsDevelopment
func (c *Setting) SetIsDevelopment(is bool) {
	defer c.lock.RUnlock()
	c.lock.RLock()
	c.isDevelopment = is
}

// GetIsDevelopment GetIsDevelopment
func (c *Setting) GetIsDevelopment() bool {
	defer c.lock.RUnlock()
	c.lock.RLock()
	return c.isDevelopment
}
