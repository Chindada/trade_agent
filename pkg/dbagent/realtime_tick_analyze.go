// Package dbagent package dbagent
package dbagent

import (
	"time"

	"gorm.io/gorm"
)

// RealTimeTickAnalyze RealTimeTickAnalyze
type RealTimeTickAnalyze struct {
	gorm.Model `json:"-" swaggerignore:"true"`

	Stock    *Stock    `json:"stock,omitempty" yaml:"stock" gorm:"foreignKey:StockID"`
	StockID  int64     `json:"stock_id,omitempty" yaml:"stock_id" gorm:"column:stock_id"`
	TickTime time.Time `json:"tick_time,omitempty" yaml:"tick_time" gorm:"column:tick_time"`

	Close      float64 `json:"close,omitempty" yaml:"close" gorm:"column:close"`
	PR         float64 `json:"pr,omitempty" yaml:"pr" gorm:"column:pr"`
	Volume     int64   `json:"volume,omitempty" yaml:"volume" gorm:"column:volume"`
	OutInRatio float64 `json:"out_in_ratio,omitempty" yaml:"out_in_ratio" gorm:"column:out_in_ratio"`
}

// TableName TableName
func (RealTimeTickAnalyze) TableName() string {
	return "realtime_tick_analyze"
}

// InsertRealTimeTickAnalyze InsertRealTimeTickAnalyze
func (c *DBAgent) InsertRealTimeTickAnalyze(record *RealTimeTickAnalyze) error {
	err := c.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&record).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}
