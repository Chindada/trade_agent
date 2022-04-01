// Package history package history
package history

import (
	"sync"
	"time"

	"trade_agent/global"
	"trade_agent/pkg/cache"
	"trade_agent/pkg/dbagent"
	"trade_agent/pkg/log"
	"trade_agent/pkg/mqhandler"
	"trade_agent/pkg/pb"
	"trade_agent/pkg/sinopacapi"
)

var wg sync.WaitGroup

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
			close, dbErr := dbagent.Get().GetHistoryCloseByStockAndDate(cache.GetCache().GetStockID(s.Stock.Number), int64(calendarDate.ID))
			if dbErr != nil {
				log.Get().Panic(err)
			}
			if close == 0 {
				var stockNumArr, dateArr []string
				dateArr = append(dateArr, t.Format(global.ShortTimeLayout))
				stockNumArr = append(stockNumArr, s.Stock.Number)

				wg.Add(1)
				err = sinopacapi.Get().FetchHistoryCloseByStockArrDateArr(stockNumArr, dateArr)
				if err != nil {
					return err
				}
				continue
			}
			cache.GetCache().SetStockHistoryClose(s.Stock.Number, close, t)
			log.Get().WithFields(map[string]interface{}{
				"Stock": s.Stock.Number,
				"Date":  t.Format(global.ShortTimeLayout),
				"Close": close,
			}).Info("History Close Already Exist")
		}
	}
	wg.Wait()
	return err
}

func stockCloseCallback(m mqhandler.MQMessage) {
	defer wg.Done()
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
			Open:         cache.GetCache().GetHistoryOpen(v.GetCode(), dateTime),
			Close:        v.GetClose(),
			Stock:        cache.GetCache().GetStock(v.GetCode()),
			CalendarDate: calendarDate,
		}
		cache.GetCache().SetStockHistoryClose(tmp.Stock.Number, tmp.Close, tmp.CalendarDate.Date)

		if err := dbagent.Get().InsertOrUpdateHistoryClose(tmp); err != nil {
			log.Get().Panic(err)
		}

		log.Get().WithFields(map[string]interface{}{
			"Stock": tmp.Stock.Number,
			"Date":  tmp.CalendarDate.Date.Format(global.ShortTimeLayout),
			"Close": tmp.Close,
		}).Info("History Close")
	}
}
