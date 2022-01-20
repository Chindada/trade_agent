// Package healthcheck package healthcheck
package healthcheck

import (
	"time"
	"trade_agent/global"
	"trade_agent/pkg/eventbus"
	"trade_agent/pkg/log"
	"trade_agent/pkg/sinopacapi"
)

var keepChannel chan struct{} = make(chan struct{})

// InitHealthCheck InitHealthCheck
func InitHealthCheck() {
	log.Get().Info("Initial HealthCheck")

	eventbus.Get().SubscribeTerminate(terminate)

	go func() {
		for range time.Tick(10 * time.Second) {
			key, err := sinopacapi.Get().FetchServerToken()
			if err != nil {
				log.Get().Error(err)
				continue
			}

			if key != sinopacapi.Get().GetToken() {
				log.Get().Error("Token expired, terminate")
				close(keepChannel)
			}
		}
	}()

	<-keepChannel
}

func terminate(t time.Time) {
	log.Get().WithFields(map[string]interface{}{
		"Time": t.Format(global.LongTimeLayout),
	}).Error("Terminate")
	close(keepChannel)
}
