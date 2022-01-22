// Package healthcheck package healthcheck
package healthcheck

import (
	"sync"
	"time"
	"trade_agent/global"
	"trade_agent/pkg/eventbus"
	"trade_agent/pkg/log"
	"trade_agent/pkg/sinopacapi"
)

var (
	keepChannel        chan struct{} = make(chan struct{})
	restartSinopacLock sync.Once
)

// InitHealthCheck InitHealthCheck
func InitHealthCheck() {
	log.Get().Info("Initial HealthCheck")

	eventbus.Get().SubscribeRestartSinopacMQSRV(restartSinopacMQSRV)

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

func restartSinopacMQSRV(t time.Time) {
	log.Get().WithFields(map[string]interface{}{
		"Time": t.Format(global.LongTimeLayout),
	}).Warn("Receive RestartSinopacSRV")

	restartSinopacLock.Do(func() {
		if err := sinopacapi.Get().RestartSinopacSRV(); err != nil {
			log.Get().Error(err)
			stuck := make(chan struct{})
			<-stuck
		}
		// ternimate self
		close(keepChannel)
	})
}
