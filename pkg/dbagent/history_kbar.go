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

// TableName TableName
func (HistoryKbar) TableName() string {
	return "history_kbar"
}

// InsertHistoryKbar InsertHistoryKbar
func (c *DBAgent) InsertHistoryKbar(record *HistoryKbar) error {
	err := c.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&record).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// InsertMultiHistoryKbar InsertMultiHistoryKbar
func (c *DBAgent) InsertMultiHistoryKbar(recordArr []*HistoryKbar) error {
	err := c.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.CreateInBatches(&recordArr, multiInsertBatchSize).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// DeleteAllHistoryKbar DeleteAllHistoryKbar
func (c *DBAgent) DeleteAllHistoryKbar() error {
	err := c.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Not("id = 0").Unscoped().Delete(&HistoryKbar{}).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// CheckHistoryKbarExistByStockID CheckHistoryKbarExistByStockID
func (c *DBAgent) CheckHistoryKbarExistByStockID(stockID int64, date time.Time) (bool, error) {
	var count int64
	if err := c.DB.Model(&HistoryKbar{}).Where("stock_id = ? and tick_time >= ? and tick_time < ?", stockID, date, date.AddDate(0, 0, 1)).Count(&count).Error; err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

// HistoryKbarArr HistoryKbarArr
type HistoryKbarArr []*HistoryKbar

// Analyzer Analyzer
func (c HistoryKbarArr) Analyzer() string {
	return ""
}

// GetHistoryKbarByStockIDAndDate GetHistoryKbarByStockIDAndDate
func (c *DBAgent) GetHistoryKbarByStockIDAndDate(stockID int64, date time.Time) (HistoryKbarArr, error) {
	var tmp HistoryKbarArr
	err := c.DB.Preload("Stock").Model(&HistoryKbar{}).Where("stock_id = ? and tick_time >= ? and tick_time < ?", stockID, date, date.AddDate(0, 0, 1)).Find(&tmp).Error
	return tmp, err
}
