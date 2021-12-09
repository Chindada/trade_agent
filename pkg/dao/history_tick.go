// Package dao package dao
package dao

import (
	"time"

	"gorm.io/gorm"
)

// HistoryTick HistoryTick
type HistoryTick struct {
	gorm.Model `json:"-" swaggerignore:"true"`
	time.Time  `json:"time" gorm:"column:time;index:idx_history_tick"`
	Stock      `json:"stock" gorm:"foreignKey:StockID"`

	StockID   int64   `json:"stock_id" gorm:"column:stock_id;index:idx_history_tick"`
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
