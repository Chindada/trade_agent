// Package cache package cache
package cache

import "fmt"

type cacheType string

const (
	calendar    cacheType = "calendar"
	targets     cacheType = "targets"
	stockDetail cacheType = "stock_detail"
	tradeDay    cacheType = "trade_day"

	biasRate     cacheType = "bias_rate"
	historyClose cacheType = "history_close"

	buyOrder       cacheType = "buyOrder"
	sellOrder      cacheType = "sellOrder"
	sellFirstOrder cacheType = "sellFirstOrder"
	buyLaterOrder  cacheType = "buyLaterOrder"
	waitingOrder   cacheType = "waitingOrder"
	forwardOrder   cacheType = "forwardOrder"
	reverseOrder   cacheType = "reverseOrder"

	realTimeBidaskChannel cacheType = "rea_time_bidask_channel"
	realTimeTickChannel   cacheType = "real_time_tick_channel"
)

func keyTypeRealTimeClose(stockNum string) cacheType {
	return cacheType(fmt.Sprintf("keyTypeRealTimeClose:%s", stockNum))
}

func keyTypeRealTimeBidAskStatus(stockNum string) cacheType {
	return cacheType(fmt.Sprintf("keyTypeRealTimeBidAskStatus:%s", stockNum))
}

func keyTypeHistoryTickAnalyze(stockNum string) cacheType {
	return cacheType(fmt.Sprintf("keyTypeHistoryTickAnalyze:%s", stockNum))
}

func keyTypeHistoryKbarAnalyze(stockNum string) cacheType {
	return cacheType(fmt.Sprintf("keyTypeHistoryKbarAnalyze:%s", stockNum))
}
