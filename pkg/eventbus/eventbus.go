// Package eventbus package eventbus
package eventbus

import (
	"sync"

	"github.com/asaskevich/EventBus"
)

var (
	globalBus EventBus.Bus
	once      sync.Once
)

func initBus() {
	if globalBus != nil {
		return
	}
	newBus := EventBus.New()
	globalBus = newBus
}

// Get Get
func Get() EventBus.Bus {
	if globalBus != nil {
		return globalBus
	}
	once.Do(initBus)
	return globalBus
}
