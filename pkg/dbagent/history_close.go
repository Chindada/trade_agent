// Package dbagent package dbagent
package dbagent

import (
	"gorm.io/gorm"
)

// HistoryClose HistoryClose
type HistoryClose struct {
	gorm.Model `json:"-" swaggerignore:"true"`

	Close          float64       `json:"close,omitempty" yaml:"close" gorm:"column:close"`
	Stock          *Stock        `json:"stock,omitempty" yaml:"stock" gorm:"foreignKey:StockID"`
	StockID        int64         `json:"stock_id,omitempty" yaml:"stock_id" gorm:"column:stock_id"`
	CalendarDate   *CalendarDate `json:"calendar_date,omitempty" yaml:"calendar_date" gorm:"foreignKey:CalendarDateID"`
	CalendarDateID int64         `json:"calendar_date_id,omitempty" yaml:"calendar_date_id" gorm:"calendar_date_id"`
}

// InsertHistoryClose InsertHistoryClose
func (c *DBAgent) InsertHistoryClose(record *HistoryClose) error {
	err := c.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&record).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// InsertMultiHistoryClose InsertMultiHistoryClose
func (c *DBAgent) InsertMultiHistoryClose(recordArr []*HistoryClose) error {
	err := c.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.CreateInBatches(&recordArr, multiInsertBatchSize).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// DeleteAllHistoryClose DeleteAllHistoryClose
func (c *DBAgent) DeleteAllHistoryClose() error {
	err := c.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Not("id = 0").Unscoped().Delete(&HistoryClose{}).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// InsertOrUpdateHistoryClose InsertOrUpdateHistoryClose
func (c *DBAgent) InsertOrUpdateHistoryClose(record *HistoryClose) error {
	err := c.DB.Transaction(func(tx *gorm.DB) error {
		var dbRecord HistoryClose
		result := tx.Model(&HistoryClose{}).
			Where("close = ?", record.Close).
			Where("stock_id = ?", record.Stock.ID).
			Where("calendar_date_id = ?", record.CalendarDate.ID).
			First(&dbRecord)
		if result.RowsAffected == 0 {
			err := c.InsertHistoryClose(record)
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

// GetHistoryCloseByStockAndDate GetHistoryCloseByStockAndDate
func (c *DBAgent) GetHistoryCloseByStockAndDate(stockID, dateID int64) (close float64, err error) {
	var tmp HistoryClose
	result := c.DB.Model(&HistoryClose{}).Where("stock_id = ? AND calendar_date_id = ?", stockID, dateID).Find(&tmp)
	return tmp.Close, result.Error
}
