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

func subStockClose(stockNumArr []string, date time.Time) {
	handler := mqhandler.Get()
	body := mqhandler.MQSubBody{
		MQTopic:  mqhandler.TopicLastcount(),
		Once:     false,
		Callback: stockCloseCallback,
	}
	err := handler.Sub(body)
	if err != nil {
		log.Get().Panic(err)
	}
	err = sinopacapi.Get().FetchStockCloseByStockArrAndDate(stockNumArr, date)
	if err != nil {
		log.Get().Panic(err)
	}
}

func stockCloseCallback(m mqhandler.MQMessage) {
	body := pb.LastCountResponse{}
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
		if err := dbagent.Get().InsertHistoryClose(tmp); err != nil {
			log.Get().Panic(err)
		}
	}
}
