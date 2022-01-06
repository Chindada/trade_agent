// Package healthchecktask package healthchecktask
package healthchecktask

import (
	"errors"
	"runtime/debug"
	"sync"
	"trade_agent/pkg/log"
	"trade_agent/pkg/sinopacapi"
)

var lock sync.Mutex

// Run Run
func Run() {
	lock.Lock()
	var err error
	defer lock.Unlock()
	defer func() {
		if r := recover(); r != nil {
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("unknown panic")
			}
			log.Get().Error(err.Error() + "\n" + string(debug.Stack()))
		}
	}()

	err = sinopacapi.Get().RestartSinopacSRV()
	if err != nil {
		log.Get().Panic(err)
	}
}
