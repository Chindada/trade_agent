// Package dbagent package dbagent
package dbagent

// Tabler Tabler
type Tabler interface {
	TableName() string
}

// TableName TableName
func (Stock) TableName() string {
	return "basic_stock"
}

// TableName TableName
func (CalendarDate) TableName() string {
	return "basic_calendar"
}

// TableName TableName
func (Target) TableName() string {
	return "trade_targets"
}

// TableName TableName
func (HistoryTick) TableName() string {
	return "history_tick"
}

// TableName TableName
func (TradeEvent) TableName() string {
	return "trade_event"
}
