// Package dbagent package dbagent
package dbagent

import (
	"gorm.io/gorm"
)

// HistoryMA HistoryMA
type HistoryMA struct {
	gorm.Model `json:"-" swaggerignore:"true"`

	QuaterMA       float64       `json:"quater_ma,omitempty" yaml:"quater_ma"`
	Stock          *Stock        `json:"stock,omitempty" yaml:"stock"`
	StockID        int64         `json:"stock_id,omitempty" yaml:"stock_id"`
	CalendarDate   *CalendarDate `json:"calendar_date,omitempty" yaml:"calendar_date"`
	CalendarDateID int64         `json:"calendar_date_id,omitempty" yaml:"calendar_date_id"`
}

// TableName TableName
func (HistoryMA) TableName() string {
	return "history_ma"
}

// InsertHistoryMA InsertHistoryMA
func (c *DBAgent) InsertHistoryMA(record *HistoryMA) error {
	err := c.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&record).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// InsertMultiHistoryMA InsertMultiHistoryMA
func (c *DBAgent) InsertMultiHistoryMA(recordArr []*HistoryMA) error {
	err := c.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.CreateInBatches(&recordArr, multiInsertBatchSize).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// DeleteAllHistoryMA DeleteAllHistoryMA
func (c *DBAgent) DeleteAllHistoryMA() error {
	err := c.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Not("id = 0").Unscoped().Delete(&HistoryMA{}).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// InsertOrUpdateHistoryMA InsertOrUpdateHistoryMA
func (c *DBAgent) InsertOrUpdateHistoryMA(record *HistoryMA) error {
	err := c.DB.Transaction(func(tx *gorm.DB) error {
		var dbRecord HistoryMA
		result := tx.Model(&HistoryMA{}).
			Where("stock_id = ?", record.Stock.ID).
			Where("calendar_date_id = ?", record.CalendarDate.ID).
			First(&dbRecord)
		if result.RowsAffected == 0 {
			err := c.InsertHistoryMA(record)
			if err != nil {
				return err
			}
		} else {
			record.Model = dbRecord.Model
			if err := tx.Save(&record).Error; err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

// GetAllQuaterMAByStockID GetAllQuaterMAByStockID
func (c *DBAgent) GetAllQuaterMAByStockID(stockID int64) (maArr []HistoryMA, err error) {
	err = c.DB.Preload("Stock").Preload("CalendarDate").Model(&HistoryMA{}).Where("stock_id = ?", stockID).Find(&maArr).Error
	return maArr, err
}

// BelowQuaterMA BelowQuaterMA
type BelowQuaterMA struct {
	Date   string  `json:"date"`
	Stocks []Stock `json:"stocks"`
}
