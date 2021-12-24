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

func subStockClose(targetArr []*dbagent.Target, fetchDate []time.Time) error {
	handler := mqhandler.Get()
	err := handler.Sub(mqhandler.MQSubBody{
		MQTopic:  mqhandler.TopicLastcountMultiDate(),
		Once:     false,
		Callback: stockCloseCallback,
	})
	if err != nil {
		return err
	}

	for _, t := range fetchDate {
		var calendarDate *dbagent.CalendarDate
		calendarDate, err = dbagent.Get().GetCalendarDate(t)
		if err != nil {
			log.Get().Panic(err)
		}
		for _, s := range targetArr {
			var close float64
			if close, err = dbagent.Get().GetHistoryCloseByStockAndDate(cache.GetCache().GetStockID(s.Stock.Number), int64(calendarDate.ID)); err != nil {
				log.Get().Panic(err)
			} else if close == 0 {
				var stockNumArr, dateArr []string
				dateArr = append(dateArr, t.Format(global.ShortTimeLayout))
				stockNumArr = append(stockNumArr, s.Stock.Number)
				err = sinopacapi.Get().FetchHistoryCloseByStockArrDateArr(stockNumArr, dateArr)
				if err != nil {
					return err
				}
				continue
			}
			cache.GetCache().Set(cache.KeyStockHistoryClose(s.Stock.Number, t), close)
			log.Get().WithFields(map[string]interface{}{
				"Stock": s.Stock.Number,
				"Date":  t.Format(global.ShortTimeLayout),
				"Close": close,
			}).Info("History Close Already Exist")
		}
	}
	return err
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
		cache.GetCache().Set(cache.KeyStockHistoryClose(tmp.Stock.Number, tmp.CalendarDate.Date), tmp.Close)

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
