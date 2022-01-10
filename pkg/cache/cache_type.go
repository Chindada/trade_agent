// Package cache package cache
package cache

import "fmt"

type keyType string

const (
	biasRate     keyType = "bias_rate"
	calendar     keyType = "calendar"
	historyClose keyType = "history_close"
	historyKbar  keyType = "history_kbar"
	historyTick  keyType = "history_tick"

	buyOrder       keyType = "buyOrder"
	sellOrder      keyType = "sellOrder"
	sellFirstOrder keyType = "sellFirstOrder"
	buyLaterOrder  keyType = "buyLaterOrder"
	waitingOrder   keyType = "waitingOrder"
	forwardOrder   keyType = "forwardOrder"
	reverseOrder   keyType = "reverseOrder"

	realTimeBidask keyType = "rea_time_bidask"
	realTimeTick   keyType = "real_time_tick"
	stockDetail    keyType = "stock_detail"
	tradeDay       keyType = "trade_day"
	targets        keyType = "targets"
)

func keyTypeRealTimeClose(stockNum string) keyType {
	return keyType(fmt.Sprintf("keyTypeRealTimeClose:%s", stockNum))
}

func keyTypeHistoryTickAnalyze(stockNum string) keyType {
	return keyType(fmt.Sprintf("keyTypeHistoryTickAnalyze:%s", stockNum))
}
