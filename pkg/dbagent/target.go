// Package dbagent package dbagent
package dbagent

import (
	"time"

	"gorm.io/gorm"
)

// Target Target
type Target struct {
	gorm.Model `json:"-" swaggerignore:"true"`

	StockID  int64     `json:"stock_id" gorm:"column:stock_id"`
	Volume   int64     `json:"volume" gorm:"column:volume"`
	Stock    Stock     `json:"stock" gorm:"foreignKey:StockID"`
	TradeDay time.Time `json:"trade_day" gorm:"column:trade_day"`
}

// InsertTarget InsertTarget
func (c *DBAgent) InsertTarget(target *Target) error {
	err := c.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&target).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// InsertMultiTarget InsertMultiTarget
func (c *DBAgent) InsertMultiTarget(targetArr []*Target) error {
	err := c.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.CreateInBatches(&targetArr, multiInsertBatchSize).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// TargetCond TargetCond
type TargetCond struct {
	LimitPriceLow  float64 `json:"limit_price_low"`
	LimitPriceHigh float64 `json:"limit_price_high"`
	LimitVolume    int64   `json:"limit_volume"`
}
