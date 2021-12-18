// Package tradeday package tradeday
package tradeday

import (
	"time"
	"trade_agent/pkg/cache"
)

// GetTradeDay GetTradeDay
func GetTradeDay() (tradeDay time.Time, err error) {
	var today time.Time
	if time.Now().Hour() >= 15 {
		today = time.Now().AddDate(0, 0, 1)
	} else {
		today = time.Now()
	}
	tradeDay, err = GetNextTradeDayTime(today)
	if err != nil {
		return tradeDay, err
	}
	return tradeDay, err
}

// GetNextTradeDayTime GetNextTradeDayTime
func GetNextTradeDayTime(nowTime time.Time) (tradeDay time.Time, err error) {
	tmp := time.Date(nowTime.Year(), nowTime.Month(), nowTime.Day(), 0, 0, 0, 0, time.Local)
	calendar := cache.GetCache().GetCalendar()
	if !calendar[tmp] {
		nowTime = nowTime.AddDate(0, 0, 1)
		return GetNextTradeDayTime(nowTime)
	}
	return tmp, err
}

// GetLastNTradeDayByDate GetLastNTradeDayByDate
func GetLastNTradeDayByDate(n int64, firstDay time.Time) []time.Time {
	calendar := cache.GetCache().GetCalendar()
	var tmp []time.Time
	for {
		if calendar[firstDay.AddDate(0, 0, -1)] {
			tmp = append(tmp, firstDay.AddDate(0, 0, -1))
		}
		if len(tmp) == int(n) {
			break
		}
		firstDay = firstDay.AddDate(0, 0, -1)
	}
	return tmp
}
