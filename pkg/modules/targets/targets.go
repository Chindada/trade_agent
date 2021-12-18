// Package targets package targets
package targets

import "trade_agent/pkg/log"

// InitTargets InitTargets
func InitTargets() {
	err := getStockTargets()
	if err != nil {
		log.Get().Panic(err)
	}
	log.Get().Info("Initial Targets")
}
