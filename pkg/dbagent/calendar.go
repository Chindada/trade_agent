// Package dbagent package dbagent
package dbagent

import (
	"time"

	"gorm.io/gorm"
)

// CalendarDate CalendarDate
type CalendarDate struct {
	gorm.Model `json:"-" swaggerignore:"true"`

	Date       time.Time `gorm:"column:date;uniqueIndex"`
	IsTradeDay bool      `gorm:"column:is_trade_day"`
}

// InsertMultiCalendarDate InsertMultiCalendarDate
func (c *DBAgent) InsertMultiCalendarDate(dateArr []*CalendarDate) error {
	err := c.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.CreateInBatches(&dateArr, multiInsertBatchSize).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// GetAllCalendarDateMap GetAllCalendarDateMap
func (c *DBAgent) GetAllCalendarDateMap() (calendarDateMap map[time.Time]bool, err error) {
	calendarDateMap = make(map[time.Time]bool)
	var calendarDateArr []CalendarDate
	result := c.DB.Model(&CalendarDate{}).Not("id = 0").Find(&calendarDateArr)
	for _, v := range calendarDateArr {
		calendarDateMap[v.Date] = true
	}
	return calendarDateMap, result.Error
}

// GetAllTradeDayMap GetAllTradeDayMap
func (c *DBAgent) GetAllTradeDayMap() (tradeDayMap map[time.Time]bool, err error) {
	tradeDayMap = make(map[time.Time]bool)
	var calendarDateArr []CalendarDate
	result := c.DB.Model(&CalendarDate{}).Where("is_trade_day = true").Find(&calendarDateArr)
	for _, v := range calendarDateArr {
		tradeDayMap[v.Date] = true
	}
	return tradeDayMap, result.Error
}
