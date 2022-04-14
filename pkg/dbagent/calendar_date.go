// Package dbagent package dbagent
package dbagent

import (
	"time"

	"gorm.io/gorm"
)

// CalendarDate CalendarDate
type CalendarDate struct {
	gorm.Model `json:"-" swaggerignore:"true"`

	Date       time.Time `json:"date,omitempty" yaml:"date" gorm:"column:date"`
	IsTradeDay bool      `json:"is_trade_day,omitempty" yaml:"is_trade_day" gorm:"column:is_trade_day"`
}

// TableName TableName
func (CalendarDate) TableName() string {
	return "basic_calendar"
}

// InsertCalendarDate InsertCalendarDate
func (c *DBAgent) InsertCalendarDate(record *CalendarDate) error {
	err := c.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&record).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// InsertMultiCalendarDate InsertMultiCalendarDate
func (c *DBAgent) InsertMultiCalendarDate(recordArr []*CalendarDate) error {
	err := c.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.CreateInBatches(&recordArr, multiInsertBatchSize).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// DeleteAllCalendarDate DeleteAllCalendarDate
func (c *DBAgent) DeleteAllCalendarDate() error {
	err := c.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Not("id = 0").Unscoped().Delete(&CalendarDate{}).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// GetAllTradeDayDate GetAllTradeDayDate
func (c *DBAgent) GetAllTradeDayDate() (arr []*CalendarDate, err error) {
	result := c.DB.Model(&CalendarDate{}).Where("is_trade_day = true").Find(&arr)
	return arr, result.Error
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
// func (c *DBAgent) GetAllTradeDayMap() (tradeDayMap map[time.Time]bool, err error) {
// 	tradeDayMap = make(map[time.Time]bool)
// 	var calendarDateArr []CalendarDate
// 	result := c.DB.Model(&CalendarDate{}).Where("is_trade_day = true").Find(&calendarDateArr)
// 	for _, v := range calendarDateArr {
// 		tradeDayMap[v.Date] = true
// 	}
// 	return tradeDayMap, result.Error
// }

// GetCalendarDate GetCalendarDate
// func (c *DBAgent) GetCalendarDate(date time.Time) (dateTime *CalendarDate, err error) {
// 	result := c.DB.Model(&CalendarDate{}).Where("date = ?", date).Find(&dateTime)
// 	return dateTime, result.Error
// }
