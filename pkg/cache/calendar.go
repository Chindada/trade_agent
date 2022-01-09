// Package cache package cache
package cache

import "time"

// KeyCalendar KeyCalendar
func KeyCalendar() *Key {
	return &Key{
		Name: "KeyCalendar",
		Type: calendar,
	}
}

// SetCalendar SetCalendar
func (c *Cache) SetCalendar(tradeDayMap map[time.Time]bool) {
	key := KeyCalendar()
	c.getCacheByType(key.Type).Set(key.Name, tradeDayMap, noExpired)
}

// GetCalendar GetCalendar
func (c *Cache) GetCalendar() map[time.Time]bool {
	c.lock.RLock()
	k := KeyCalendar()
	tmp := c.CacheMap[string(k.Type)]
	c.lock.RUnlock()
	if tmp == nil {
		return nil
	}
	if value, ok := tmp.Get(k.Name); ok {
		return value.(map[time.Time]bool)
	}
	return nil
}
