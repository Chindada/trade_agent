// Package tickprocess package tickprocess
package tickprocess

import (
	"trade_agent/pkg/cache"
	"trade_agent/pkg/dbagent"
	"trade_agent/pkg/log"
	"trade_agent/pkg/mqhandler"
	"trade_agent/pkg/pb"
)

func subHistroyTick() error {
	handler := mqhandler.Get()
	err := handler.Sub(mqhandler.MQSubBody{
		MQTopic:  mqhandler.TopicHistoryTick(),
		Once:     false,
		Callback: historyTickCallback,
	})
	if err != nil {
		return err
	}
	return nil
}

func historyTickCallback(m mqhandler.MQMessage) {
	body := pb.HistoryTickResponse{}
	err := body.UnmarshalProto(m.Payload())
	if err != nil {
		log.Get().Panic(err)
	}

	var saveTick dbagent.HistoryTickArr
	for _, v := range body.GetData() {
		saveTick = append(saveTick, v.ToHistoryTick(body.GetStockNum()))
	}

	if err := dbagent.Get().InsertMultiHistoryTick(saveTick); err != nil {
		log.Get().Panic(err)
	}

	// analyze to cache
	analyzeVolumeArr := saveTick.Analyzer()
	cache.GetCache().Set(cache.KeyStockHistoryTickAnalyze(body.GetStockNum()), analyzeVolumeArr)

	log.Get().WithFields(map[string]interface{}{
		"Stock": body.GetStockNum(),
		"Date":  body.GetDate(),
	}).Info("HistoryTick Process Done")
}
