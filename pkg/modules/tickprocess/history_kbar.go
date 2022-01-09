// Package tickprocess package tickprocess
package tickprocess

import (
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

	if err := dbagent.Get().InsertMultiHistoryKbar(saveKbar); err != nil {
		log.Get().Panic(err)
	}

	if kbarStatus := saveKbar.Analyzer(); kbarStatus != "" {
		cache.GetCache().SetStockHistoryKbarAnalyze(body.GetStockNum(), kbarStatus)
	}

	log.Get().WithFields(map[string]interface{}{
		"Stock": body.GetStockNum(),
		"Date":  body.GetStartDate(),
	}).Info("HistoryKbar Process Done")
}
