// Package eventbus package eventbus
package eventbus

import (
	"github.com/asaskevich/EventBus"
)

var globalBus *EventBus.Bus

func newBus() *EventBus.Bus {
	bus := EventBus.New()
	return &bus
}

// Get Get
func Get() *EventBus.Bus {
	if globalBus != nil {
		return globalBus
	}
	return newBus()
}
