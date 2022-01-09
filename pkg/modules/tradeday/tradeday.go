// Package tradeday package tradeday
package tradeday

import (
	"time"
	"trade_agent/global"
	"trade_agent/pkg/cache"
	"trade_agent/pkg/config"
	"trade_agent/pkg/dbagent"
	"trade_agent/pkg/log"
)

// InitTradeDay InitTradeDay
func InitTradeDay() {
	log.Get().Info("Initial TradeDay")

	// save calendar to db and cache
	err := ImportCalendar()
	if err != nil {
		log.Get().Panic(err)
	}

	tradeDayMap, err := dbagent.Get().GetAllTradeDayMap()
	if err != nil {
		log.Get().Panic(err)
	}
	cache.GetCache().SetCalendar(tradeDayMap)

	// get trade day
	tradeDay, err := GetTradeDay()
	if err != nil {
		log.Get().Panic(err)
	}
	// save to cache
	cache.GetCache().SetTradeDay(tradeDay)
	log.Get().WithFields(map[string]interface{}{
		"Date": tradeDay.Format(global.ShortTimeLayout),
	}).Info("TradeDay")

	// every 10 seconds to check if now is open time
	tradeConf := config.GetTradeConfig()
	go func() {
		for range time.Tick(10 * time.Second) {
			isOpen := checkIsOpenTimeWithEndWaitTime(tradeDay, tradeConf.TradeInEndTime, tradeConf.WaitInOpen)
			cache.GetCache().SetIsOpenWithEndWaitTime(isOpen)
		}
	}()

	// get config
	closeRange := GetLastNTradeDayByDate(config.GetTradeConfig().HistoryClosePeriod, tradeDay)
	tickRange := GetLastNTradeDayByDate(config.GetTradeConfig().HistoryTickPeriod, tradeDay)
	kbarRange := GetLastNTradeDayByDate(config.GetTradeConfig().HistoryKbarPeriod, tradeDay)

	// save to cahce
	cache.GetCache().SetHistroyCloseRange(closeRange)
	cache.GetCache().SetHistroyTickRange(tickRange)
	cache.GetCache().SetHistroyKbarRange(kbarRange)
}

func checkIsOpenTimeWithEndWaitTime(tradeDay time.Time, tradInEndTime, waitInOpen int64) bool {
	starTime := tradeDay.Add(9*time.Hour + time.Duration(waitInOpen)*time.Minute)
	if time.Now().After(starTime) && time.Now().Before(starTime.Add(time.Duration(tradInEndTime)*time.Hour)) {
		return true
	}
	return false
}
