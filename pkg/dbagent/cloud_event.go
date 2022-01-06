// Package dbagent package dbagent
package dbagent

import (
	"gorm.io/gorm"
)

// CloudEvent CloudEvent
type CloudEvent struct {
	gorm.Model `json:"-" swaggerignore:"true"`

	Event     string `json:"event,omitempty" yaml:"event" gorm:"column:event"`
	EventCode int64  `json:"event_code,omitempty" yaml:"event_code" gorm:"column:event_code"`
	Info      string `json:"info,omitempty" yaml:"info" gorm:"column:info"`
	Response  int64  `json:"response,omitempty" yaml:"response" gorm:"column:response"`
}

// TableName TableName
func (CloudEvent) TableName() string {
	return "cloud_event"
}

// InsertCloudEvent InsertCloudEvent
func (c *DBAgent) InsertCloudEvent(record *CloudEvent) error {
	err := c.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&record).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// InsertMultiCloudEvent InsertMultiCloudEvent
func (c *DBAgent) InsertMultiCloudEvent(recordArr []*CloudEvent) error {
	err := c.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.CreateInBatches(&recordArr, multiInsertBatchSize).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// DeleteAllCloudEvent DeleteAllCloudEvent
func (c *DBAgent) DeleteAllCloudEvent() error {
	err := c.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Not("id = 0").Delete(&CloudEvent{}).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}
