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

	"google.golang.org/protobuf/proto"
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
	if err := proto.Unmarshal(m.Payload(), &body); err != nil {
		log.Get().Errorf("Format Wrong: %s", string(m.Payload()))
		return
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
		if err := dbagent.Get().InsertOrUpdateHistoryClose(tmp); err != nil {
			log.Get().Panic(err)
		}
	}
}
