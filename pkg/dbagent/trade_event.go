// Package dbagent package dbagent
package dbagent

import (
	"gorm.io/gorm"
)

// TradeEvent TradeEvent
type TradeEvent struct {
	gorm.Model `json:"-" swaggerignore:"true"`

	Event     string `json:"event,omitempty" yaml:"event" gorm:"column:event"`
	EventCode int64  `json:"event_code,omitempty" yaml:"event_code" gorm:"column:event_code"`
	Info      string `json:"info,omitempty" yaml:"info" gorm:"column:info"`
	Response  int64  `json:"response,omitempty" yaml:"response" gorm:"column:response"`
}

// InsertTradeEvent InsertTradeEvent
func (c *DBAgent) InsertTradeEvent(record *TradeEvent) error {
	err := c.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&record).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// InsertMultiTradeEvent InsertMultiTradeEvent
func (c *DBAgent) InsertMultiTradeEvent(recordArr []*TradeEvent) error {
	err := c.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.CreateInBatches(&recordArr, multiInsertBatchSize).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// DeleteAllTradeEvent DeleteAllTradeEvent
func (c *DBAgent) DeleteAllTradeEvent() error {
	err := c.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Not("id = 0").Unscoped().Delete(&TradeEvent{}).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}
