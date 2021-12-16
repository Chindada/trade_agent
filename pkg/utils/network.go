// Package utils package utils
package utils

import (
	"net"
	"sync"
	"time"
	"trade_agent/pkg/log"
)

// lock for check port
var checkLock sync.Mutex

// CheckPortIsOpen CheckPortIsOpen
func CheckPortIsOpen(host string, port string) bool {
	defer checkLock.Unlock()
	checkLock.Lock()

	timeout := time.Second
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), timeout)
	if err != nil {
		log.Get().Warn(err)
	}
	if conn != nil {
		defer func() {
			if err := conn.Close(); err != nil {
				log.Get().Error(err)
			}
		}()
		return true
	}
	return false
}

// GetHostIP GetHostIP
func GetHostIP() string {
	var results []string
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Get().Error(err)
	}
	var addrs []net.Addr
	for _, i := range ifaces {
		if i.HardwareAddr.String() == "" {
			continue
		}
		addrs, err = i.Addrs()
		if err != nil {
			log.Get().Error(err)
		}
		for _, addr := range addrs {
			if ip := addr.(*net.IPNet).IP.To4(); ip != nil {
				if ip[0] != 127 && ip[0] != 169 {
					results = append(results, ip.String())
				}
			}
		}
	}
	return results[len(results)-1]
}
