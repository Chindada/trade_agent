// Package tickprocess package tickprocess
package tickprocess

import (
	"trade_agent/pkg/cache"
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

	// skip if simtrade
	if body.GetTick().GetSimtrade() == 1 {
		simTradeRealTimeTickChannel <- 1
		return
	}

	// adapter to local struct
	tick := body.GetTick().ToRealTimeTick()

	// send to channel
	if ch := cache.GetCache().GetRealTimeTickChannel(tick.Stock.Number); ch != nil {
		ch <- tick
	}

	// save to db
	// if err := dbagent.Get().InsertRealTimeTick(tick); err != nil {
	// 	log.Get().Panic(err)
	// }
}
