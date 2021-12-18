// Package order package order
package order

import "trade_agent/pkg/log"

// InitOrder InitOrder
func InitOrder() {
	err := updateOrderStatus()
	if err != nil {
		log.Get().Panic(err)
	}
	log.Get().Info("Initial Order")
}
