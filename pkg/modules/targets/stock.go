// Package targets package targets
package targets

import (
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

// getStockTargets getStockTargets
func getStockTargets() error {
	wg.Add(1)
	handler := mqhandler.Get()
	err := handler.Sub(mqhandler.MQSubBody{
		MQTopic:  mqhandler.TopicSnapshotAll(),
		Once:     true,
		Callback: snapshotAllCallback,
	})
	if err != nil {
		return err
	}

	err = sinopacapi.Get().FetchAllSnapShot()
	if err != nil {
		return err
	}
	wg.Wait()
	return nil
}

// snapshotAllCallback snapshotAllCallback
func snapshotAllCallback(m mqhandler.MQMessage) {
	defer wg.Done()
	body := pb.SnapshotResponse{}
	if err := proto.Unmarshal(m.Payload(), &body); err != nil {
		log.Get().Errorf("Format Wrong: %s", string(m.Payload()))
		return
	}

	conf, err := config.Get()
	if err != nil {
		log.Get().Panic(err)
	}
	condition := conf.GetTradeTargetCondtion()

	var targetArr []*dbagent.Target
	for _, v := range body.GetData() {
		if stockTargetFilter(v, condition) {
			tmp := &dbagent.Target{
				Stock:    cache.GetCache().GetStock(v.GetCode()),
				TradeDay: cache.GetCache().GetTradeDay(),
				Volume:   v.GetTotalVolume(),
			}
			targetArr = append(targetArr, tmp)
			log.Get().WithFields(map[string]interface{}{
				"Stock":       tmp.Stock.Name,
				"TotalVolume": tmp.Volume,
			}).Info("Target")
		}
	}

	err = dbagent.Get().DeleteMultiTargetByDate(cache.GetCache().GetTradeDay())
	if err != nil {
		log.Get().Panic(err)
	}
	err = dbagent.Get().InsertMultiTarget(targetArr)
	if err != nil {
		log.Get().Panic(err)
	}

	// send to bus
	eventbus.Get().Pub(eventbus.TopicTargets(), targetArr)
}

func stockTargetFilter(v *pb.SnapshotMessage, cond *config.TargetCond) bool {
	if v.GetTotalVolume() < cond.LimitVolume {
		return false
	}
	if v.GetClose() < cond.LimitPriceLow || v.GetClose() > cond.LimitPriceHigh {
		return false
	}
	if !cache.GetCache().GetStock(v.GetCode()).DayTrade {
		return false
	}
	return true
}
