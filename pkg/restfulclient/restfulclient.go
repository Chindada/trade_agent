// Package restfulclient package restfulclient
package restfulclient

import (
	"sync"
	"trade_agent/pkg/log"

	"github.com/go-resty/resty/v2"
)

var (
	globalClient *resty.Client
	once         sync.Once
)

func initClient() {
	if globalClient != nil {
		return
	}
	globalClient = resty.New()
	globalClient.SetLogger(log.Get())
}

// Get Get
func Get() *resty.Client {
	if globalClient != nil {
		return globalClient
	}
	once.Do(initClient)
	return globalClient
}
