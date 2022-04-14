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
	var stockIDArr []uint
	targetMap := make(map[uint]*dbagent.Target)
	for _, v := range targetArr {
		stockIDArr = append(stockIDArr, v.Stock.ID)
		targetMap[v.Stock.ID] = v
	}
	noDataDateArr := []tempFetch{}
	for _, t := range fetchDate {
		calendarDate := cache.GetCache().GetCalendarID(t.Format(global.ShortTimeLayout))
		closeMap, dbErr := dbagent.Get().GetHistoryCloseByMultiStockAndDate(stockIDArr, int64(calendarDate.ID))
		if dbErr != nil {
			log.Get().Panic(err)
		}
		for s, v := range closeMap {
			if v == 0 {
				noDataDateArr = append(noDataDateArr, tempFetch{
					target: targetMap[s],
					date:   t,
				})
			} else {
				cache.GetCache().SetStockHistoryClose(targetMap[s].Stock.Number, v, t)
				log.Get().WithFields(map[string]interface{}{
					"Stock": targetMap[s].Stock.Number,
					"Date":  t.Format(global.ShortTimeLayout),
					"Close": v,
				}).Info("History Close Already Exist")
			}
		}
	}
	for _, v := range noDataDateArr {
		var stockNumArr, dateArr []string
		dateArr = append(dateArr, v.date.Format(global.ShortTimeLayout))
		stockNumArr = append(stockNumArr, v.target.Stock.Number)

		wg.Add(1)
		err = sinopacapi.Get().FetchHistoryCloseByStockArrDateArr(stockNumArr, dateArr)
		if err != nil {
			return err
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

		calendarDate := cache.GetCache().GetCalendarID(dateTime.Format(global.ShortTimeLayout))
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
