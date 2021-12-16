// Package tradeday package tradeday
package tradeday

import (
	"time"
	"trade_agent/global"
	"trade_agent/pkg/cache"
	"trade_agent/pkg/dbagent"
	"trade_agent/pkg/log"
)

// InitTradeDay InitTradeDay
func InitTradeDay() {
	// save calendar to db and cache
	err := ImportCalendar()
	if err != nil {
		log.Get().Panic(err)
	}
	tradeDayMap, err := dbagent.Get().GetAllTradeDayMap()
	if err != nil {
		log.Get().Panic(err)
	}
	cache.GetCache().Set(cache.KeyCalendar(), tradeDayMap)

	// get trade day
	tradeDay, err := GetTradeDay()
	if err != nil {
		log.Get().Panic(err)
	}
	cache.GetCache().Set(cache.KeyTradeDay(), tradeDay)

	log.Get().WithFields(map[string]interface{}{
		"Date": tradeDay.Format(global.ShortTimeLayout),
	}).Info("TradeDay")
}

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
func GetLastNTradeDayByDate(n int, firstDay time.Time) []time.Time {
	calendar := cache.GetCache().GetCalendar()
	var tmp []time.Time
	for {
		if calendar[firstDay.AddDate(0, 0, -1)] {
			tmp = append(tmp, firstDay.AddDate(0, 0, -1))
		}
		if len(tmp) == n {
			break
		}
		firstDay = firstDay.AddDate(0, 0, -1)
	}
	return tmp
}
