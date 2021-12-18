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
	return "basic_targets"
}

// TableName TableName
func (CloudEvent) TableName() string {
	return "cloud_event"
}

// TableName TableName
func (HistoryTick) TableName() string {
	return "history_tick"
}

// TableName TableName
func (Balance) TableName() string {
	return "trade_balance"
}

// TableName TableName
func (HistoryBidAsk) TableName() string {
	return "history_bidask"
}

// TableName TableName
func (HistoryKbar) TableName() string {
	return "history_kbar"
}

// TableName TableName
func (RealTimeTick) TableName() string {
	return "realtime_tick"
}

// TableName TableName
func (OrderStatus) TableName() string {
	return "order_status"
}

// TableName TableName
func (AnalyzedTick) TableName() string {
	return "analyzed_tick"
}

// TableName TableName
func (SimulationResult) TableName() string {
	return "simulation_result"
}

// TableName TableName
func (SimulationCondition) TableName() string {
	return "simulation_condition"
}

// TableName TableName
func (HistoryClose) TableName() string {
	return "history_close"
}
