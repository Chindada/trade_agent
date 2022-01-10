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

const (
	openTime time.Duration = 9 * time.Hour
)

// InitTradeDay InitTradeDay
func InitTradeDay() {
	log.Get().Info("Initial TradeDay")

	// save calendar to db and cache
	err := ImportCalendar()
	if err != nil {
		log.Get().Panic(err)
	}

	// get trade config
	tradeConf := config.GetTradeConfig()

	// save trade day map to cache
	tradeDayMap, err := dbagent.Get().GetAllTradeDayMap()
	if err != nil {
		log.Get().Panic(err)
	}
	cache.GetCache().SetCalendar(tradeDayMap)

	// get trade day and save to cache
	tradeDay, err := GetTradeDay()
	if err != nil {
		log.Get().Panic(err)
	}
	cache.GetCache().SetTradeDay(tradeDay)

	// trade out time
	tradeOutEndTime := tradeDay.Add(openTime).Add(time.Duration(tradeConf.TradeOutEndTime) * time.Hour)
	cache.GetCache().SetTradeDayTradeOutEndTime(tradeOutEndTime)

	log.Get().WithFields(map[string]interface{}{
		"Date":            tradeDay.Format(global.ShortTimeLayout),
		"TradeOutEndTime": tradeOutEndTime.Format(global.LongTimeLayout),
	}).Info("TradeDay")

	// every 10 seconds to check if now is open time
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
	starTime := tradeDay.Add(openTime + time.Duration(waitInOpen)*time.Minute)
	if time.Now().After(starTime) && time.Now().Before(starTime.Add(time.Duration(tradInEndTime)*time.Hour)) {
		return true
	}
	return false
}
