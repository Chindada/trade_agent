// Package targets package targets
package targets

import (
	"sync"
	"trade_agent/pkg/cache"
	"trade_agent/pkg/log"
	"trade_agent/pkg/mqhandler"
	"trade_agent/pkg/pb"
	"trade_agent/pkg/sinopacapi"
)

var tseWg sync.WaitGroup

// getTSERealTime getTSERealTime
func getTSERealTime() error {
	handler := mqhandler.Get()
	err := handler.Sub(mqhandler.MQSubBody{
		MQTopic:  mqhandler.TopicSnapshotTSE(),
		Once:     false,
		Callback: tseSnapShotCallback,
	})
	if err != nil {
		return err
	}

	tseWg.Add(1)
	go func() {
		defer tseWg.Done()
		err = sinopacapi.Get().FetchTSESnapShot()
		if err != nil {
			log.Get().Error(err)
			return
		}
	}()
	tseWg.Wait()
	return nil
}

// tseSnapShotCallback tseSnapShotCallback
func tseSnapShotCallback(m mqhandler.MQMessage) {
	body := pb.SnapshotResponse{}
	err := body.UnmarshalProto(m.Payload())
	if err != nil {
		log.Get().Panic(err)
	}

	for _, v := range body.GetData() {
		stock := cache.GetCache().GetStock(v.GetCode())
		if stock == nil {
			log.Get().WithFields(map[string]interface{}{
				"Stock": v.GetCode(),
			}).Error("Stock Cache Error")
			continue
		}
		cache.GetCache().SetTSESnapshot(v.ToTSESnapshot())
	}
}
