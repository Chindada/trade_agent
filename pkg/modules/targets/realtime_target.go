// Package targets package targets
package targets

import (
	"sort"
	"sync"
	"trade_agent/pkg/cache"
	"trade_agent/pkg/config"
	"trade_agent/pkg/dbagent"
	"trade_agent/pkg/eventbus"
	"trade_agent/pkg/log"
	"trade_agent/pkg/mqhandler"
	"trade_agent/pkg/pb"
	"trade_agent/pkg/sinopacapi"
)

var wg sync.WaitGroup

// getRealTimeTargets getRealTimeTargets
func getRealTimeTargets() error {
	handler := mqhandler.Get()
	err := handler.Sub(mqhandler.MQSubBody{
		MQTopic:  mqhandler.TopicSnapshotAll(),
		Once:     false,
		Callback: snapShotCallback,
	})
	if err != nil {
		return err
	}
	wg.Add(1)
	err = sinopacapi.Get().FetchAllSnapShot()
	if err != nil {
		return err
	}
	wg.Wait()
	return nil
}

// snapShotCallback snapShotCallback
func snapShotCallback(m mqhandler.MQMessage) {
	defer wg.Done()
	body := pb.SnapshotResponse{}
	err := body.UnmarshalProto(m.Payload())
	if err != nil {
		log.Get().Panic(err)
	}
	sort.Slice(body.GetData(), func(i, j int) bool {
		return body.GetData()[i].GetTotalVolume() > body.GetData()[j].GetTotalVolume()
	})
	condition := config.GetTargetCondConfig()
	tradeDay := cache.GetCache().GetTradeDay()
	var targetArr []*dbagent.Target
	var newTargetCount int
	for _, v := range body.GetData() {
		stock := cache.GetCache().GetStock(v.GetCode())
		if stock == nil {
			log.Get().WithFields(map[string]interface{}{
				"Stock": v.GetCode(),
			}).Error("Stock Cache Error")
			continue
		}
		tmpTarget := stockWithData{
			stock:       stock,
			close:       v.GetClose(),
			totalVolume: v.GetTotalVolume(),
		}
		if stockTargetFilter(tmpTarget, condition, true) {
			newTargetCount++
			if exist, dbErr := dbagent.Get().CheckExistTargetsByDateStockID(tradeDay, int64(tmpTarget.stock.ID)); dbErr != nil {
				log.Get().Error(err)
				continue
			} else if !exist {
				tmp := &dbagent.Target{
					Stock:       tmpTarget.stock,
					TradeDay:    tradeDay,
					Volume:      tmpTarget.totalVolume,
					Rank:        len(targetArr) + 1 + 100,
					RealTimeAdd: true,
				}
				targetArr = append(targetArr, tmp)
				log.Get().WithFields(map[string]interface{}{
					"Stock":       tmp.Stock.Name,
					"TotalVolume": tmp.Volume,
					"Rank":        tmp.Rank,
				}).Info("RealTime Target")
			}
			if newTargetCount >= int(condition.RealTimeTargetsCount) {
				break
			}
		}
	}
	if len(targetArr) == 0 {
		return
	} else if err := dbagent.Get().InsertMultiTarget(targetArr); err != nil {
		log.Get().Panic(err)
	}
	// send to bus
	eventbus.Get().Pub(eventbus.TopicTargets(), targetArr)
	// append to cache
	cache.GetCache().AppendTargets(targetArr)
}
