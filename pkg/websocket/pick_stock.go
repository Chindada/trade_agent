package websocket

import (
	"sync"
	"time"

	"trade_agent/pkg/cache"
	"trade_agent/pkg/log"
	"trade_agent/pkg/mqhandler"
	"trade_agent/pkg/pb"
	"trade_agent/pkg/sinopacapi"
)

var (
	processPickStockLock sync.Mutex
	wg                   sync.WaitGroup
	snapShotMap          = make(map[string]*pb.SnapshotMessage)
	snapShotMapLock      sync.RWMutex
)

// getRealTimeTargets getRealTimeTargets
func getRealTimeTargets(stockNumArr []string) error {
	handler := mqhandler.Get()
	err := handler.Sub(mqhandler.MQSubBody{
		MQTopic:  mqhandler.TopicSnapshots(),
		Once:     false,
		Callback: fetchSnapShotsCallback,
	})
	if err != nil {
		log.Get().Error(err)
		return err
	}

	wg.Add(1)
	var tmp []string
	for _, v := range stockNumArr {
		if cache.GetCache().GetStock(v) != nil {
			tmp = append(tmp, v)
		}
	}
	if len(tmp) == 0 {
		wg.Done()
		return nil
	}
	err = sinopacapi.Get().FetchSnapShots(tmp)
	if err != nil {
		log.Get().Error(err)
		return err
	}
	wg.Wait()
	return nil
}

// fetchSnapShotsCallback fetchSnapShotsCallback
func fetchSnapShotsCallback(m mqhandler.MQMessage) {
	defer wg.Done()
	body := pb.SnapshotResponse{}
	err := body.UnmarshalProto(m.Payload())
	if err != nil {
		log.Get().Panic(err)
	}
	snapShotMapLock.Lock()
	for _, v := range body.GetData() {
		snapShotMap[v.GetCode()] = v
	}
	snapShotMapLock.Unlock()
}

func processPickStock() {
	defer processPickStockLock.Unlock()
	processPickStockLock.Lock()
	for {
		time.Sleep(1 * time.Second)
		userMapLock.Lock()
		liveClients := activeUser
		userMapLock.Unlock()

		var allStock []string
		tmpStockMap := make(map[string]struct{})
		for _, w := range liveClients {
			for _, v := range w.PickStock {
				if _, ok := tmpStockMap[v]; !ok {
					allStock = append(allStock, v)
				}
			}
		}
		if len(allStock) == 0 {
			continue
		}

		if err := getRealTimeTargets(allStock); err != nil {
			log.Get().Error(err)
		}

		for _, w := range liveClients {
			if len(w.PickStock) == 0 {
				continue
			}
			w.lock.Lock()
			w.Msg = getRealTimeStock(w.PickStock)
			w.lock.Unlock()
			go func(m *wsClient) {
				msgChan <- m
			}(w)
		}
	}
}

// SocketPickStock SocketPickStock
type SocketPickStock struct {
	StockNum        string  `json:"stock_num" yaml:"stock_num"`
	StockName       string  `json:"stock_name" yaml:"stock_name"`
	IsTarget        bool    `json:"is_target" yaml:"is_target"`
	PriceChange     float64 `json:"price_change" yaml:"price_change"`
	PriceChangeRate float64 `json:"price_change_rate" yaml:"price_change_rate"`
	Price           float64 `json:"price" yaml:"price"`
	Wrong           bool    `json:"wrong" yaml:"wrong"`
}

func getRealTimeStock(stockList []string) []SocketPickStock {
	snapShotMapLock.RLock()
	dataMap := snapShotMap
	snapShotMapLock.RUnlock()
	var result []SocketPickStock

	for _, v := range stockList {
		var name string
		if s := cache.GetCache().GetStock(v); s != nil {
			name = s.Name
		} else {
			result = append(result, SocketPickStock{
				StockNum: v,
				Wrong:    true,
			})
			continue
		}
		result = append(result, SocketPickStock{
			StockNum:        v,
			StockName:       name,
			IsTarget:        false,
			PriceChange:     dataMap[v].GetChangePrice(),
			PriceChangeRate: dataMap[v].GetChangeRate(),
			Price:           dataMap[v].GetClose(),
		})
	}
	return result
}
