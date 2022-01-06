// Package tickprocess package tickprocess
package tickprocess

import (
	"trade_agent/pkg/cache"
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

	// skip if simtrade
	if body.GetBidAsk().GetSimtrade() == 1 {
		return
	}

	// adapter to local struct
	bidAsk := body.GetBidAsk().ToRealTimeBidAsk()

	// send to channel
	if ch := cache.GetCache().GetRealTimeBidAskChannel(bidAsk.Stock.Number); ch != nil {
		ch <- bidAsk
	}

	// save to db
	err = dbagent.Get().InsertRealTimeBidAsk(bidAsk)
	if err != nil {
		log.Get().Panic(err)
	}
}
