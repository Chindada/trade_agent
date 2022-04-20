// Package tickprocess package tickprocess
package tickprocess

import (
	"time"

	"trade_agent/global"
	"trade_agent/pkg/cache"
	"trade_agent/pkg/dbagent"
	"trade_agent/pkg/log"
	"trade_agent/pkg/mqhandler"
	"trade_agent/pkg/pb"
)

func subHistroyKbar() error {
	handler := mqhandler.Get()
	err := handler.Sub(mqhandler.MQSubBody{
		MQTopic:  mqhandler.TopicHistoryKbar(),
		Once:     false,
		Callback: historyKbarCallback,
	})
	if err != nil {
		return err
	}
	return nil
}

func historyKbarCallback(m mqhandler.MQMessage) {
	body := pb.HistoryKbarResponse{}
	err := body.UnmarshalProto(m.Payload())
	if err != nil {
		log.Get().Panic(err)
	}

	var saveKbar dbagent.HistoryKbarArr
	for _, v := range body.GetData() {
		saveKbar = append(saveKbar, v.ToHistoryKbar(body.GetStockNum()))
	}

	if len(saveKbar) == 0 {
		return
	}

	dataTime, err := time.ParseInLocation(global.ShortTimeLayout, body.GetStartDate(), time.Local)
	if err != nil {
		log.Get().Panic(err)
	}
	cache.GetCache().SetStockHistoryOpen(body.GetStockNum(), saveKbar[0].Open, dataTime)
	var close, open, high, low float64
	var volume int64
	var lastTickTime time.Time
	for i, kbar := range saveKbar {
		if i == 0 {
			open = kbar.Open
		}
		if i == len(saveKbar)-1 {
			close = kbar.Close
			lastTickTime = kbar.TickTime
		}
		if high == 0 {
			high = kbar.High
		} else if kbar.High > high {
			high = kbar.High
		}
		if low == 0 {
			low = kbar.Low
		} else if kbar.Low < low {
			low = kbar.Low
		}
		volume += kbar.Volume
	}
	cache.GetCache().SetStockHistoryDayKbar(body.GetStartDate(), saveKbar[0].Stock.Number, &dbagent.HistoryKbar{
		Close:    close,
		Open:     open,
		High:     high,
		Low:      low,
		Volume:   volume,
		TickTime: lastTickTime,
	})

	if err := dbagent.Get().InsertMultiHistoryKbar(saveKbar); err != nil {
		log.Get().Panic(err)
	}

	// if kbarStatus := saveKbar.Analyzer(); kbarStatus != "" {
	// 	cache.GetCache().SetStockHistoryKbarAnalyze(body.GetStockNum(), kbarStatus)
	// }

	log.Get().WithFields(map[string]interface{}{
		"Stock": body.GetStockNum(),
		"Date":  body.GetStartDate(),
	}).Info("HistoryKbar Process Done")
}
