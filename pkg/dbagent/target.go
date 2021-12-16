// Package dbagent package dbagent
package dbagent

import (
	"time"

	"gorm.io/gorm"
)

// TargetCond TargetCond
type TargetCond struct {
	LimitPriceLow  float64 `json:"limit_price_low,omitempty" yaml:"limit_price_low" gorm:"column:limit_price_low"`
	LimitPriceHigh float64 `json:"limit_price_high,omitempty" yaml:"limit_price_high" gorm:"column:limit_price_high"`
	LimitVolume    int64   `json:"limit_volume,omitempty" yaml:"limit_volume" gorm:"column:limit_volume"`
}

// Target Target
type Target struct {
	gorm.Model `json:"-" swaggerignore:"true"`

	Stock    *Stock    `json:"stock,omitempty" yaml:"stock" gorm:"foreignKey:StockID"`
	StockID  int64     `json:"stock_id,omitempty" yaml:"stock_id" gorm:"column:stock_id"`
	TradeDay time.Time `json:"trade_day,omitempty" yaml:"trade_day" gorm:"column:trade_day"`

	Volume int64 `json:"volume,omitempty" yaml:"volume" gorm:"column:volume"`
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

// DeleteMultiTargetByDate DeleteMultiTargetByDate
func (c *DBAgent) DeleteMultiTargetByDate(date time.Time) error {
	err := c.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("trade_day = ?", date).Unscoped().Delete(&Target{}).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}
