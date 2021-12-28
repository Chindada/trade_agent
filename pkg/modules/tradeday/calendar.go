// Package tradeday package tradeday
package tradeday

import (
	"fmt"
	"time"
	"trade_agent/global"
	"trade_agent/pkg/dbagent"
	"trade_agent/pkg/log"
)

const (
	startTradeYear int = 2021
	endTradeYear   int = 2022
)

var holidayArr = []string{
	// 2021
	"2021-01-01", "2021-02-08", "2021-02-09",
	"2021-02-10", "2021-02-11", "2021-02-12",
	"2021-02-15", "2021-02-16", "2021-03-01",
	"2021-04-02", "2021-04-05", "2021-04-30",
	"2021-06-14", "2021-09-20", "2021-09-21",
	"2021-10-11", "2021-12-31",

	// 2022
	"2022-01-27", "2022-01-28", "2022-01-31",
	"2022-02-01", "2022-02-02", "2022-02-03",
	"2022-02-04", "2022-02-05", "2022-02-06",
	"2022-02-28", "2022-04-04", "2022-04-05",
	"2022-05-02", "2022-06-03", "2022-09-09",
	"2022-10-10",
}

const (
	zeroTime string = "00:00:00"
)

// ImportCalendar ImportCalendar
func ImportCalendar() (err error) {
	holidayTimeMap := make(map[time.Time]bool)
	for _, h := range holidayArr {
		var holidayTime time.Time
		holidayString := fmt.Sprintf("%s %s", h, zeroTime)
		holidayTime, err = time.ParseInLocation(global.LongTimeLayout, holidayString, time.Local)
		if err != nil {
			return err
		}
		holidayTimeMap[holidayTime] = true
	}
	inDBCalendarMap, err := dbagent.Get().GetAllCalendarDateMap()
	if err != nil {
		return err
	}
	firstDay := time.Date(startTradeYear, 1, 1, 0, 0, 0, 0, time.Local)
	var calendarDateArr []*dbagent.CalendarDate
	var exist, insert int
	for {
		var isTradeDay bool
		if firstDay.Year() > endTradeYear {
			break
		}
		if firstDay.Weekday() != time.Saturday && firstDay.Weekday() != time.Sunday && !holidayTimeMap[firstDay] {
			isTradeDay = true
		}
		tmp := &dbagent.CalendarDate{
			Date:       firstDay,
			IsTradeDay: isTradeDay,
		}
		if !inDBCalendarMap[firstDay] {
			calendarDateArr = append(calendarDateArr, tmp)
			insert++
		} else {
			exist++
		}
		firstDay = firstDay.AddDate(0, 0, 1)
	}
	err = dbagent.Get().InsertMultiCalendarDate(calendarDateArr)
	log.Get().WithFields(map[string]interface{}{
		"Exist":  exist,
		"Insert": insert,
	}).Info("ImportCalendar")
	return err
}
