// Package dbagent package dbagent
package dbagent

import "gorm.io/gorm"

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
