// Package tickprocess package tickprocess
package tickprocess

import (
	"trade_agent/pkg/dbagent"
	"trade_agent/pkg/log"
	"trade_agent/pkg/mqhandler"
	"trade_agent/pkg/pb"

	"google.golang.org/protobuf/proto"
)

func subRealTimeTick() error {
	handler := mqhandler.Get()
	// realtime
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
	if err := proto.Unmarshal(m.Payload(), &body); err != nil {
		log.Get().Errorf("Format Wrong: %s", string(m.Payload()))
		return
	}
	if err := dbagent.Get().InsertRealTimeTick(body.GetTick().ToRealTimeTick()); err != nil {
		log.Get().Panic(err)
	}
}
