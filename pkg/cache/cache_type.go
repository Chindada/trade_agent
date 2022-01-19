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

	buyOrder       cacheType = "buy_order"
	sellOrder      cacheType = "sell_order"
	sellFirstOrder cacheType = "sellFirst_order"
	buyLaterOrder  cacheType = "buyLater_order"
	waitingOrder   cacheType = "waiting_order"
	forwardOrder   cacheType = "forward_order"
	reverseOrder   cacheType = "reverse_order"

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
