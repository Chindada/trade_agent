// Package tradeday package tradeday
package tradeday

import (
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
	cache.GetCache().Set(cache.KeyCalendar(), tradeDayMap)

	// get trade day
	tradeDay, err := GetTradeDay()
	if err != nil {
		log.Get().Panic(err)
	}
	// save to cache
	cache.GetCache().Set(cache.KeyTradeDay(), tradeDay)
	log.Get().WithFields(map[string]interface{}{
		"Date": tradeDay.Format(global.ShortTimeLayout),
	}).Info("TradeDay")

	closeRange := GetLastNTradeDayByDate(config.GetTradeConfig().HistoryClosePeriod, tradeDay)
	tickRange := GetLastNTradeDayByDate(config.GetTradeConfig().HistoryTickPeriod, tradeDay)
	kbarRange := GetLastNTradeDayByDate(config.GetTradeConfig().HistoryKbarPeriod, tradeDay)

	cache.GetCache().Set(cache.KeyHistroyCloseRange(), closeRange)
	cache.GetCache().Set(cache.KeyHistroyTickRange(), tickRange)
	cache.GetCache().Set(cache.KeyHistroyKbarRange(), kbarRange)
}
