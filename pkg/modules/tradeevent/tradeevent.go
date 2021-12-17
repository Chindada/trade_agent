// Package tradeevent package tradeevent
package tradeevent

import (
	"trade_agent/pkg/dbagent"
	"trade_agent/pkg/log"
	"trade_agent/pkg/mqhandler"
	"trade_agent/pkg/pb"

	"google.golang.org/protobuf/proto"
)

// InitTradeEvent InitTradeEvent
func InitTradeEvent() {
	handler := mqhandler.Get()
	body := mqhandler.MQSubBody{
		MQTopic:  mqhandler.TopicTradeEvent(),
		Once:     false,
		Callback: tredeEventCallback,
	}
	err := handler.Sub(body)
	if err != nil {
		log.Get().Panic(err)
	}
}

func tredeEventCallback(m mqhandler.MQMessage) {
	var err error
	body := pb.EventResponse{}
	if err = proto.Unmarshal(m.Payload(), &body); err != nil {
		log.Get().Errorf("Format Wrong: %s", string(m.Payload()))
		return
	}
	err = dbagent.Get().InsertTradeEvent(body.ToTradeEvent())
	if err != nil {
		log.Get().Panic(err)
	}
	log.Get().WithFields(map[string]interface{}{
		"EventCode": body.GetEventCode(),
		"Event":     body.GetEvent(),
		"RespCode":  body.GetRespCode(),
		"Info":      body.GetInfo(),
	}).Info("TradeEvent")
}
