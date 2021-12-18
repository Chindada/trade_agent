// Package tickprocess package tickprocess
package tickprocess

import (
	"trade_agent/pkg/dbagent"
	"trade_agent/pkg/log"
	"trade_agent/pkg/mqhandler"
	"trade_agent/pkg/pb"

	"google.golang.org/protobuf/proto"
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
	if err := proto.Unmarshal(m.Payload(), &body); err != nil {
		log.Get().Errorf("Format Wrong: %s", string(m.Payload()))
		return
	}
	var saveTick []*dbagent.HistoryTick
	for _, v := range body.GetData() {
		saveTick = append(saveTick, v.ToHistoryTick(body.GetStockNum()))
	}
	if err := dbagent.Get().InsertMultiHistoryTick(saveTick); err != nil {
		log.Get().Panic(err)
	}
}
