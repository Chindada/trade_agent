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
