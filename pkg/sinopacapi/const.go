// Package sinopacapi package sinopacapi
package sinopacapi

const (
	// StatusSuccuss StatusSuccuss
	StatusSuccuss string = "success"
	// StatusFail StatusFail
	StatusFail string = "fail"
	// StatusCancelOrderNotFound StatusCancelOrderNotFound
	StatusCancelOrderNotFound string = "cancel order not found"
	// StatusAlreadyCanceled StatusAlreadyCanceled
	StatusAlreadyCanceled string = "order already be canceled"
)

// OrderAction OrderAction
type OrderAction int64

const (
	// ActionBuy ActionBuy
	ActionBuy OrderAction = iota + 1
	// ActionSell ActionSell
	ActionSell
	// ActionSellFirst ActionSellFirst
	ActionSellFirst
)

// TickType TickType
type TickType int64

const (
	// TickTypeStockRealTime TickTypeStockRealTime
	TickTypeStockRealTime TickType = iota + 1
	// TickTypeStockBidAsk TickTypeStockBidAsk
	TickTypeStockBidAsk
)
