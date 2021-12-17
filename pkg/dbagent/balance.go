// Package dbagent package dbagent
package dbagent

import (
	"time"

	"gorm.io/gorm"
)

// Balance Balance
type Balance struct {
	gorm.Model `json:"-" swaggerignore:"true"`

	TradeDay time.Time `json:"trade_day,omitempty" yaml:"trade_day" gorm:"column:trade_day"`

	TradeCount      int64 `json:"trade_count,omitempty" yaml:"trade_count" gorm:"column:trade_count"`
	Forward         int64 `json:"forward,omitempty" yaml:"forward" gorm:"column:forward"`
	Reverse         int64 `json:"reverse,omitempty" yaml:"reverse" gorm:"column:reverse"`
	OriginalBalance int64 `json:"original_balance,omitempty" yaml:"original_balance" gorm:"column:original_balance"`
	Discount        int64 `json:"discount,omitempty" yaml:"discount" gorm:"column:discount"`
	Total           int64 `json:"total,omitempty" yaml:"total" gorm:"column:total"`
}

// InsertBalance InsertBalance
func (c *DBAgent) InsertBalance(record *Balance) error {
	err := c.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&record).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// InsertMultiBalance InsertMultiBalance
func (c *DBAgent) InsertMultiBalance(recordArr []*Balance) error {
	err := c.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.CreateInBatches(&recordArr, multiInsertBatchSize).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// DeleteAllBalance DeleteAllBalance
func (c *DBAgent) DeleteAllBalance() error {
	err := c.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Not("id = 0").Unscoped().Delete(&Balance{}).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// InsertOrUpdateBalance InsertOrUpdateBalance
func (c *DBAgent) InsertOrUpdateBalance(record *Balance) error {
	err := c.DB.Transaction(func(tx *gorm.DB) error {
		var inDB Balance
		dbResult := c.DB.Model(&Balance{}).Where("trade_day = ?", record.TradeDay).Find(&inDB)
		if dbResult.Error != nil {
			return dbResult.Error
		}
		if dbResult.RowsAffected == 0 {
			if err := tx.Create(&record).Error; err != nil {
				return err
			}
		} else {
			record.Model = inDB.Model
			if err := tx.Save(&record).Error; err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

// DeleteMultiBalanceByDate DeleteMultiBalanceByDate
func (c *DBAgent) DeleteMultiBalanceByDate(date time.Time) error {
	err := c.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("trade_day = ?", date).Unscoped().Delete(&Balance{}).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}
