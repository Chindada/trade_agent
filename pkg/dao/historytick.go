// Package dao package dao
package dao

import (
	"time"

	"gorm.io/gorm"
)

// HistoryTick HistoryTick
type HistoryTick struct {
	gorm.Model `json:"-" swaggerignore:"true"`
	TickTime   time.Time `json:"tick_time" gorm:"column:tick_time;index:idx_history_tick"`

	StockID int64 `json:"stock_id" gorm:"column:stock_id;index:idx_history_tick"`
	Stock   Stock `json:"stock" gorm:"foreignKey:StockID"`

	Close     float64 `json:"close" gorm:"column:close"`
	TickType  int64   `json:"tick_type" gorm:"column:tick_type"`
	Volume    int64   `json:"volume" gorm:"column:volume"`
	BidPrice  float64 `json:"bid_price" gorm:"column:bid_price"`
	BidVolume int64   `json:"bid_volume" gorm:"column:bid_volume"`
	AskPrice  float64 `json:"ask_price" gorm:"column:ask_price"`
	AskVolume int64   `json:"ask_volume" gorm:"column:ask_volume"`

	Open float64 `json:"open" gorm:"column:open"`
	High float64 `json:"high" gorm:"column:high"`
	Low  float64 `json:"low" gorm:"column:low"`
}
