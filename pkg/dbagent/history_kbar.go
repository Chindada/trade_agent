// Package dbagent package dbagent
package dbagent

import (
	"time"

	"gorm.io/gorm"
)

// HistoryKbar HistoryKbar
type HistoryKbar struct {
	gorm.Model `json:"-" swaggerignore:"true"`

	Stock    *Stock    `json:"stock,omitempty" yaml:"stock" gorm:"foreignKey:StockID"`
	StockID  int64     `json:"stock_id,omitempty" yaml:"stock_id" gorm:"column:stock_id"`
	TickTime time.Time `json:"tick_time,omitempty" yaml:"tick_time" gorm:"column:tick_time"`

	Close  float64 `json:"close,omitempty" yaml:"close" gorm:"column:close"`
	Open   float64 `json:"open,omitempty" yaml:"open" gorm:"column:open"`
	High   float64 `json:"high,omitempty" yaml:"high" gorm:"column:high"`
	Low    float64 `json:"low,omitempty" yaml:"low" gorm:"column:low"`
	Volume int64   `json:"volume,omitempty" yaml:"volume" gorm:"column:volume"`
}
