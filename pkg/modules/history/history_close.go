// Package history package history
package history

import (
	"time"
	"trade_agent/global"
	"trade_agent/pkg/cache"
	"trade_agent/pkg/dbagent"
	"trade_agent/pkg/log"
	"trade_agent/pkg/mqhandler"
	"trade_agent/pkg/pb"
	"trade_agent/pkg/sinopacapi"
)

func subStockClose(targetArr []*dbagent.Target, date []time.Time) {
	handler := mqhandler.Get()
	err := handler.Sub(mqhandler.MQSubBody{
		MQTopic:  mqhandler.TopicLastcountMultiDate(),
		Once:     false,
		Callback: stockCloseCallback,
	})
	if err != nil {
		log.Get().Panic(err)
	}

	var stockNumArr, dateArr []string
	for _, t := range date {
		dateArr = append(dateArr, t.Format(global.ShortTimeLayout))
	}
	for _, s := range targetArr {
		stockNumArr = append(stockNumArr, s.Stock.Number)
	}

	err = sinopacapi.Get().FetchHistoryCloseByStockArrDateArr(stockNumArr, dateArr)
	if err != nil {
		log.Get().Panic(err)
	}
}

func stockCloseCallback(m mqhandler.MQMessage) {
	body := pb.HistoryCloseResponse{}
	err := body.UnmarshalProto(m.Payload())
	if err != nil {
		log.Get().Panic(err)
	}

	for _, v := range body.GetData() {
		dateTime, err := time.ParseInLocation(global.ShortTimeLayout, v.GetDate(), time.Local)
		if err != nil {
			log.Get().Panic(err)
		}
		calendarDate, err := dbagent.Get().GetCalendarDate(dateTime)
		if err != nil {
			log.Get().Panic(err)
		}
		tmp := &dbagent.HistoryClose{
			Close:        v.GetClose(),
			Stock:        cache.GetCache().GetStock(v.GetCode()),
			CalendarDate: calendarDate,
		}
		log.Get().WithFields(map[string]interface{}{
			"Stock": tmp.Stock.Number,
			"Date":  tmp.CalendarDate.Date.Format(global.ShortTimeLayout),
			"Close": tmp.Close,
		}).Info("History Close")
		if err := dbagent.Get().InsertOrUpdateHistoryClose(tmp); err != nil {
			log.Get().Panic(err)
		}
	}
}
