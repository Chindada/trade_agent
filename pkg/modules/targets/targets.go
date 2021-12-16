// Package targets package targets
package targets

import (
	"encoding/json"
	"sync"
	"trade_agent/pkg/cache"
	"trade_agent/pkg/config"
	"trade_agent/pkg/dbagent"
	"trade_agent/pkg/eventbus"
	"trade_agent/pkg/log"
	"trade_agent/pkg/mqhandler"
	"trade_agent/pkg/pb"
	"trade_agent/pkg/sinopacapi"

	"google.golang.org/protobuf/proto"
)

var wg sync.WaitGroup

// InitTargets InitTargets
func InitTargets() {
	handler := mqhandler.Get()
	wg.Add(1)
	body := mqhandler.MQSubBody{
		MQTopic:  mqhandler.TopicSnapshotAll(),
		Once:     true,
		Callback: snapShotAllCallback,
	}
	err := handler.Sub(body)
	if err != nil {
		log.Get().Panic(err)
	}
	err = sinopacapi.Get().FetchAllSnapShot()
	if err != nil {
		log.Get().Panic(err)
	}
	wg.Wait()
}

func snapShotAllCallback(m mqhandler.MQMessage) {
	defer wg.Done()
	var err error
	body := pb.SnapshotResponse{}
	if err = proto.Unmarshal(m.Payload(), &body); err != nil {
		log.Get().Errorf("Format Wrong: %s", string(m.Payload()))
		return
	}
	conf, err := config.Get()
	if err != nil {
		log.Get().Panic(err)
	}

	tradeConf := conf.GetTradeConfig()
	var cond dbagent.TargetCond
	err = json.Unmarshal([]byte(tradeConf.TargetCondition), &cond)
	if err != nil {
		log.Get().Panic(err)
	}

	var targetArr []*dbagent.Target
	for _, v := range body.GetData() {
		if targetFilter(v, cond) {
			stock := cache.GetCache().GetStock(v.GetCode())
			tmp := &dbagent.Target{
				Stock:    stock,
				TradeDay: cache.GetCache().GetTradeDay(),
				Volume:   v.GetTotalVolume(),
			}
			targetArr = append(targetArr, tmp)
		}
	}
	// eventbus.Get().Publish(eventbus.TopicTargets(), targetArr)
	eventbus.Get().Pub(eventbus.TopicTargets(), targetArr)

	err = dbagent.Get().DeleteMultiTargetByDate(cache.GetCache().GetTradeDay())
	if err != nil {
		log.Get().Panic(err)
	}
	err = dbagent.Get().InsertMultiTarget(targetArr)
	if err != nil {
		log.Get().Panic(err)
	}
}

func targetFilter(v *pb.SnapshotMessage, cond dbagent.TargetCond) bool {
	if v.GetTotalVolume() < cond.LimitVolume {
		return false
	}
	if v.GetClose() < cond.LimitPriceLow || v.GetClose() > cond.LimitPriceHigh {
		return false
	}
	return true
}
