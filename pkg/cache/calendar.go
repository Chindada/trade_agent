// Package cache package cache
package cache

import "time"

// KeyCalendar KeyCalendar
func KeyCalendar() string {
	return "Calendar"
}

// GetCalendar GetCalendar
func (c *Cache) GetCalendar() map[time.Time]bool {
	defer c.lock.RUnlock()
	c.lock.RLock()
	if value, ok := c.Cache.Get(KeyCalendar()); ok {
		return value.(map[time.Time]bool)
	}
	return nil
}
