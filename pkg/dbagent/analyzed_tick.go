// Package dbagent package dbagent
package dbagent

import (
	"time"

	"gorm.io/gorm"
)

// AnalyzedTick AnalyzedTick
type AnalyzedTick struct {
	gorm.Model `json:"-" swaggerignore:"true"`

	Stock    *Stock    `json:"stock,omitempty" yaml:"stock" gorm:"foreignKey:StockID"`
	StockID  int64     `json:"stock_id,omitempty" yaml:"stock_id" gorm:"column:stock_id"`
	TickTime time.Time `json:"tick_time,omitempty" yaml:"tick_time" gorm:"column:tick_time"`

	Close            float64 `json:"close,omitempty" yaml:"close" gorm:"column:close"`
	CloseChangeRatio float64 `json:"close_change_ratio,omitempty" yaml:"close_change_ratio" gorm:"column:close_change_ratio"`
	OpenChangeRatio  float64 `json:"open_change_ratio,omitempty" yaml:"open_change_ratio" gorm:"column:open_change_ratio"`
	OutSum           int64   `json:"out_sum,omitempty" yaml:"out_sum" gorm:"column:out_sum"`
	InSum            int64   `json:"in_sum,omitempty" yaml:"in_sum" gorm:"column:in_sum"`
	Volume           int64   `json:"volume,omitempty" yaml:"volume" gorm:"column:volume"`
	TotalTime        float64 `json:"total_time,omitempty" yaml:"total_time" gorm:"column:total_time"`
	CloseDiff        float64 `json:"close_diff,omitempty" yaml:"close_diff" gorm:"column:close_diff"`
	Open             float64 `json:"open,omitempty" yaml:"open" gorm:"column:open"`
	High             float64 `json:"high,omitempty" yaml:"high" gorm:"column:high"`
	Low              float64 `json:"low,omitempty" yaml:"low" gorm:"column:low"`
}

// InsertAnalyzedTick InsertAnalyzedTick
func (c *DBAgent) InsertAnalyzedTick(record *AnalyzedTick) error {
	err := c.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&record).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// InsertMultiAnalyzedTick InsertMultiAnalyzedTick
func (c *DBAgent) InsertMultiAnalyzedTick(recordArr []*AnalyzedTick) error {
	err := c.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.CreateInBatches(&recordArr, multiInsertBatchSize).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// DeleteAllAnalyzedTick DeleteAllAnalyzedTick
func (c *DBAgent) DeleteAllAnalyzedTick() error {
	err := c.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Not("id = 0").Unscoped().Delete(&AnalyzedTick{}).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}
