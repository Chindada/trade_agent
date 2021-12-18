// Package cloudevent package cloudevent
package cloudevent

import (
	"trade_agent/pkg/dbagent"
	"trade_agent/pkg/log"
	"trade_agent/pkg/mqhandler"
	"trade_agent/pkg/pb"

	"google.golang.org/protobuf/proto"
)

// InitCloudEvent InitCloudEvent
func InitCloudEvent() {
	err := updateTradeEvent()
	if err != nil {
		log.Get().Panic(err)
	}
	log.Get().Info("Initial CloudEvent")
}

// updateTradeEvent updateTradeEvent
func updateTradeEvent() error {
	handler := mqhandler.Get()
	err := handler.Sub(mqhandler.MQSubBody{
		MQTopic:  mqhandler.TopicTradeEvent(),
		Once:     false,
		Callback: tredeEventCallback,
	})
	if err != nil {
		return err
	}
	return nil
}

// tredeEventCallback tredeEventCallback
func tredeEventCallback(m mqhandler.MQMessage) {
	var err error
	body := pb.EventResponse{}
	if err = proto.Unmarshal(m.Payload(), &body); err != nil {
		log.Get().Errorf("Format Wrong: %s", string(m.Payload()))
		return
	}
	err = dbagent.Get().InsertCloudEvent(body.ToTradeEvent())
	if err != nil {
		log.Get().Panic(err)
	}
	if body.GetEventCode() != 16 {
		log.Get().WithFields(map[string]interface{}{
			"EventCode": body.GetEventCode(),
			"Event":     body.GetEvent(),
			"RespCode":  body.GetRespCode(),
			"Info":      body.GetInfo(),
		}).Warn("TradeEvent")
	}
}
