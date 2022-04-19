// Package eventbus package eventbus
package eventbus

import (
	"sync"
	"time"

	"trade_agent/pkg/dbagent"
	"trade_agent/pkg/eventbus/bustopic"
	"trade_agent/pkg/log"
	"trade_agent/pkg/sinopacapi"

	"github.com/asaskevich/EventBus"
)

// BusAgent BusAgent
type BusAgent struct {
	bus EventBus.Bus
}

var (
	globalBus *BusAgent
	once      sync.Once
)

func initBus() {
	if globalBus != nil {
		return
	}
	newAgent := &BusAgent{
		bus: EventBus.New(),
	}
	globalBus = newAgent
}

// Get Get
func Get() *BusAgent {
	if globalBus != nil {
		return globalBus
	}
	once.Do(initBus)
	return globalBus
}

// PublishTargets PublishTargets
func (c *BusAgent) PublishTargets(targetArr []*dbagent.Target) {
	go c.bus.Publish(bustopic.Targets, targetArr)
}

// SubscribeTargets SubscribeTargets
func (c *BusAgent) SubscribeTargets(f func(targetArr []*dbagent.Target)) {
	err := c.bus.Subscribe(bustopic.Targets, f)
	if err != nil {
		log.Get().Panic(err)
	}
}

// PublishNeedAnalyzeTargets PublishNeedAnalyzeTargets
func (c *BusAgent) PublishNeedAnalyzeTargets(targetArr []*dbagent.Target) {
	go c.bus.Publish(bustopic.NeedAnalyzeTargets, targetArr)
}

// SubscribeNeedAnalyzeTargets SubscribeNeedAnalyzeTargets
func (c *BusAgent) SubscribeNeedAnalyzeTargets(f func(targetArr []*dbagent.Target)) {
	err := c.bus.Subscribe(bustopic.NeedAnalyzeTargets, f)
	if err != nil {
		log.Get().Panic(err)
	}
}

// PublishStockOrder PublishStockOrder
func (c *BusAgent) PublishStockOrder(order *sinopacapi.Order) {
	go c.bus.Publish(bustopic.StockOrder, order)
}

// SubscribeStockOrder SubscribeStockOrder
func (c *BusAgent) SubscribeStockOrder(f func(order *sinopacapi.Order)) {
	err := c.bus.Subscribe(bustopic.StockOrder, f)
	if err != nil {
		log.Get().Panic(err)
	}
}

// PublishSubscribeTargets PublishSubscribeTargets
func (c *BusAgent) PublishSubscribeTargets(targetArr []*dbagent.Target) {
	go c.bus.Publish(bustopic.SubscribeTargets, targetArr)
}

// SubscribeSubscribeTargets SubscribeSubscribeTargets
func (c *BusAgent) SubscribeSubscribeTargets(f func(targetArr []*dbagent.Target)) {
	err := c.bus.Subscribe(bustopic.SubscribeTargets, f)
	if err != nil {
		log.Get().Panic(err)
	}
}

// PublishUnSubscribeTargets PublishUnSubscribeTargets
func (c *BusAgent) PublishUnSubscribeTargets(target *dbagent.Target) {
	go c.bus.Publish(bustopic.UnSubscribeTargets, target)
}

// SubscribeUnSubscribeTargets SubscribeUnSubscribeTargets
func (c *BusAgent) SubscribeUnSubscribeTargets(f func(target *dbagent.Target)) {
	err := c.bus.Subscribe(bustopic.UnSubscribeTargets, f)
	if err != nil {
		log.Get().Panic(err)
	}
}

// PublishRestartSinopacMQSRV PublishRestartSinopacMQSRV
func (c *BusAgent) PublishRestartSinopacMQSRV(t time.Time) {
	go c.bus.Publish(bustopic.RestartSinopacMQSRV, t)
}

// SubscribeRestartSinopacMQSRV SubscribeRestartSinopacMQSRV
func (c *BusAgent) SubscribeRestartSinopacMQSRV(f func(t time.Time)) {
	err := c.bus.Subscribe(bustopic.RestartSinopacMQSRV, f)
	if err != nil {
		log.Get().Panic(err)
	}
}
