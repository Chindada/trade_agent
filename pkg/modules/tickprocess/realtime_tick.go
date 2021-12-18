// Package tickprocess package tickprocess
package tickprocess

import (
	"trade_agent/pkg/dbagent"
	"trade_agent/pkg/log"
	"trade_agent/pkg/mqhandler"
	"trade_agent/pkg/pb"
)

func subRealTimeTick() error {
	handler := mqhandler.Get()
	err := handler.Sub(mqhandler.MQSubBody{
		MQTopic:  mqhandler.TopicRealTimeTick(),
		Once:     false,
		Callback: realTimeTickCallback,
	})
	if err != nil {
		return err
	}
	return nil
}

func realTimeTickCallback(m mqhandler.MQMessage) {
	body := pb.RealTimeTickResponse{}
	err := body.UnmarshalProto(m.Payload())
	if err != nil {
		log.Get().Panic(err)
	}

	if err := dbagent.Get().InsertRealTimeTick(body.GetTick().ToRealTimeTick()); err != nil {
		log.Get().Panic(err)
	}
}
