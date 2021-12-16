// Package eventbus package eventbus
package eventbus

import (
	"sync"

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

// Pub Pub
func (c *BusAgent) Pub(topic string, arg interface{}) {
	go c.bus.Publish(topic, arg)
}

// Sub Sub
func (c *BusAgent) Sub(topic string, fn interface{}) error {
	return c.bus.Subscribe(topic, fn)
}
