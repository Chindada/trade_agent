// Package dbagent package dbagent
package dbagent

import (
	"time"

	"gorm.io/gorm"
)

// Target Target
type Target struct {
	gorm.Model `json:"-" swaggerignore:"true"`

	Stock    *Stock    `json:"stock,omitempty" yaml:"stock" gorm:"foreignKey:StockID"`
	StockID  int64     `json:"stock_id,omitempty" yaml:"stock_id" gorm:"column:stock_id"`
	TradeDay time.Time `json:"trade_day,omitempty" yaml:"trade_day" gorm:"column:trade_day"`

	Rank   int   `json:"rank,omitempty" yaml:"rank" gorm:"column:rank"`
	Volume int64 `json:"volume,omitempty" yaml:"volume" gorm:"column:volume"`
}

// InsertTarget InsertTarget
func (c *DBAgent) InsertTarget(record *Target) error {
	err := c.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&record).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// InsertMultiTarget InsertMultiTarget
func (c *DBAgent) InsertMultiTarget(recordArr []*Target) error {
	err := c.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.CreateInBatches(&recordArr, multiInsertBatchSize).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// DeleteAllTarget DeleteAllTarget
func (c *DBAgent) DeleteAllTarget() error {
	err := c.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Not("id = 0").Unscoped().Delete(&Target{}).Error; err != nil {
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

// GetTargetsByDate GetTargetsByDate
func (c *DBAgent) GetTargetsByDate(date time.Time) (targetArr []*Target, err error) {
	result := c.DB.Preload("Stock").Where("trade_day = ?", date).Find(&targetArr)
	return targetArr, result.Error
}
