// Package order package order
package order

import (
	"time"
	"trade_agent/pkg/dbagent"
	"trade_agent/pkg/log"
	"trade_agent/pkg/mqhandler"
	"trade_agent/pkg/pb"
	"trade_agent/pkg/sinopacapi"

	"google.golang.org/protobuf/proto"
)

func updateOrderStatus() error {
	handler := mqhandler.Get()
	err := handler.Sub(mqhandler.MQSubBody{
		MQTopic:  mqhandler.TopicOrderStatus(),
		Once:     false,
		Callback: orderStausCallback,
	})
	if err != nil {
		return err
	}
	go func() {
		for range time.Tick(1*time.Second + 500*time.Millisecond) {
			if err := sinopacapi.Get().FetchOrderStatus(); err != nil {
				log.Get().Error(err)
			}
		}
	}()
	return nil
}

func orderStausCallback(m mqhandler.MQMessage) {
	body := pb.OrderStatusHistoryResponse{}
	if err := proto.Unmarshal(m.Payload(), &body); err != nil {
		log.Get().Errorf("Format Wrong: %s", string(m.Payload()))
		return
	}
	var saveStatus []*dbagent.OrderStatus
	for _, v := range body.GetData() {
		saveStatus = append(saveStatus, v.ToOrderStatus())
	}
	err := dbagent.Get().InsertOrUpdateMultiOrderStatus(saveStatus)
	if err != nil {
		log.Get().Error(err)
	}
}
