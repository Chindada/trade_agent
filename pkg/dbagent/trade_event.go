// Package dbagent package dbagent
package dbagent

import (
	"gorm.io/gorm"
)

// TradeEvent TradeEvent
type TradeEvent struct {
	gorm.Model `json:"-" swaggerignore:"true"`

	Event     string `gorm:"column:event" json:"event"`
	EventCode int64  `gorm:"column:event_code" json:"event_code"`
	Info      string `gorm:"column:info" json:"info"`
	Response  int64  `gorm:"column:response" json:"response"`
}

// InsertTradeEvent InsertTradeEvent
func (c *DBAgent) InsertTradeEvent(event TradeEvent) error {
	err := c.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&event).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// DeleteAllTradeEvent DeleteAllTradeEvent
func (c *DBAgent) DeleteAllTradeEvent() error {
	err := c.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Not("id = 0").Delete(&TradeEvent{}).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}
