// Package tickprocess package tickprocess
package tickprocess

import (
	"trade_agent/pkg/dbagent"
	"trade_agent/pkg/log"
	"trade_agent/pkg/mqhandler"
	"trade_agent/pkg/pb"
)

func subRealTimeBidAsk() error {
	handler := mqhandler.Get()
	err := handler.Sub(mqhandler.MQSubBody{
		MQTopic:  mqhandler.TopicRealTimeBidask(),
		Once:     false,
		Callback: realTimeBidAskCallback,
	})
	if err != nil {
		return err
	}
	return nil
}

func realTimeBidAskCallback(m mqhandler.MQMessage) {
	body := pb.RealTimeBidAskResponse{}
	err := body.UnmarshalProto(m.Payload())
	if err != nil {
		log.Get().Panic(err)
	}
	err = dbagent.Get().InsertRealTimeBidAsk(body.GetBidAsk().ToRealTimeBidAsk())
	if err != nil {
		log.Get().Panic(err)
	}
}
