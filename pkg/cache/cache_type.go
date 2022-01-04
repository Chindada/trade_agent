// Package cache package cache
package cache

// keyType keyType
type keyType string

const (
	biasRate       keyType = "bias_rate"
	calendar       keyType = "calendar"
	historyClose   keyType = "history_close"
	historyKbar    keyType = "history_kbar"
	historyTick    keyType = "history_tick"
	order          keyType = "order"
	realTimeBidask keyType = "rea_time_bidask"
	realTimeTick   keyType = "real_time_tick"
	stockDetail    keyType = "stock_detail"
	tradeDay       keyType = "trade_day"
)
