// Package restfulclient package restfulclient
package restfulclient

import (
	"github.com/go-resty/resty/v2"
	"gitlab.tocraw.com/root/toc_trader/pkg/log"
)

var globalClient *resty.Client

// Get Get
func Get() *resty.Client {
	if globalClient != nil {
		return globalClient
	}
	new := resty.New()
	new.SetLogger(log.Get())
	globalClient = new
	return globalClient
}
