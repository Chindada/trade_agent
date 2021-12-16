// Package tickprocess package tickprocess
package tickprocess

import (
	"trade_agent/pkg/dbagent"
	"trade_agent/pkg/log"
	"trade_agent/pkg/mqhandler"
	"trade_agent/pkg/pb"

	"google.golang.org/protobuf/proto"
)

func realTimeTickCallback(m mqhandler.MQMessage) {
	body := pb.RealTimeTickResponse{}
	if err := proto.Unmarshal(m.Payload(), &body); err != nil {
		log.Get().Errorf("Format Wrong: %s", string(m.Payload()))
		return
	}
	tick := body.GetTick()
	if err := dbagent.Get().InsertRealTimeTick(tick.ToRealTimeTick()); err != nil {
		log.Get().Panic(err)
	}
}
