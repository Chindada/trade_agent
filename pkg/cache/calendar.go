// Package cache package cache
package cache

import (
	"fmt"
	"time"

	"trade_agent/global"
	"trade_agent/pkg/dbagent"
)

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
	k := KeyCalendar()
	if value, ok := c.getCacheByType(k.Type).Get(k.Name); ok {
		return value.(map[time.Time]bool)
	}
	return nil
}

// KeyCalendarID KeyCalendarID
func KeyCalendarID(date string) *Key {
	return &Key{
		Name: fmt.Sprintf("KeyCalendarID:%s", date),
		Type: calendar,
	}
}

// SetCalendarID SetCalendarID
func (c *Cache) SetCalendarID(dateTime *dbagent.CalendarDate) {
	key := KeyCalendarID(dateTime.Date.Format(global.ShortTimeLayout))
	c.getCacheByType(key.Type).Set(key.Name, dateTime, noExpired)
}

// GetCalendarID GetCalendarID
func (c *Cache) GetCalendarID(date string) *dbagent.CalendarDate {
	k := KeyCalendarID(date)
	if value, ok := c.getCacheByType(k.Type).Get(k.Name); ok {
		return value.(*dbagent.CalendarDate)
	}
	return nil
}
