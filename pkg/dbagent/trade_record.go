// Package dbagent package dbagent
package dbagent

import (
	"time"

	"gorm.io/gorm"
)

// TradeRecord TradeRecord
type TradeRecord struct {
	gorm.Model `json:"-" swaggerignore:"true"`

	Stock    *Stock    `json:"stock,omitempty" yaml:"stock" gorm:"foreignKey:StockID"`
	StockID  int64     `json:"stock_id,omitempty" yaml:"stock_id" gorm:"column:stock_id"`
	TickTime time.Time `json:"tick_time,omitempty" yaml:"tick_time" gorm:"column:tick_time"`

	Action    int64     `json:"action,omitempty" yaml:"action" gorm:"column:action"`
	Price     float64   `json:"price,omitempty" yaml:"price" gorm:"column:price"`
	Quantity  int64     `json:"quantity,omitempty" yaml:"quantity" gorm:"column:quantity"`
	Status    int64     `json:"status,omitempty" yaml:"status" gorm:"column:status"`
	OrderID   string    `json:"order_id,omitempty" yaml:"order_id" gorm:"column:order_id"`
	OrderTime time.Time `json:"order_time,omitempty" yaml:"order_time" gorm:"column:order_time"`
	BuyCost   int64     `json:"buy_cost,omitempty" yaml:"buy_cost" gorm:"column:buy_cost"`
	TradeTime time.Time `json:"trade_time,omitempty" yaml:"trade_time" gorm:"column:trade_time"`
}
