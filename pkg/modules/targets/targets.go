// Package targets package targets
package targets

import (
	"encoding/json"
	"sync"
	"time"
	"trade_agent/pkg/cache"
	"trade_agent/pkg/config"
	"trade_agent/pkg/dbagent"
	"trade_agent/pkg/eventbus"
	"trade_agent/pkg/log"
	"trade_agent/pkg/mqhandler"
	"trade_agent/pkg/pb"
	"trade_agent/pkg/sinopacapi"

	mqtt "github.com/eclipse/paho.mqtt.golang"
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
		Callback: snapShotAllCallback(),
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

func snapShotAllCallback() mqhandler.MQCallback {
	return func(c mqtt.Client, m mqtt.Message) {
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
			if v.GetTotalVolume() < cond.LimitVolume {
				continue
			}
			if v.GetClose() < cond.LimitPriceLow || v.GetClose() > cond.LimitPriceHigh {
				continue
			}
			stock := cache.GetCache().Get(cache.KeyStockDetail(v.Code)).(*dbagent.Stock)
			tmp := &dbagent.Target{
				Stock:    *stock,
				TradeDay: cache.GetCache().Get(cache.KeyTradeDay()).(time.Time),
				Volume:   v.GetTotalVolume(),
			}
			targetArr = append(targetArr, tmp)
			eventbus.Get().Publish(eventbus.BusTopicTargets(), tmp)
		}
		err = dbagent.Get().InsertMultiTarget(targetArr)
		if err != nil {
			log.Get().Panic(err)
		}
	}
}
