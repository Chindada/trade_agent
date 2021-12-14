// Package eventbus package eventbus
package eventbus

import "fmt"

// BusTopicStockDetail BusTopicStockDetail
func BusTopicStockDetail(stockNum string) string {
	return fmt.Sprintf("StockDetail:%s", stockNum)
}

// BusTopicTargets BusTopicTargets
func BusTopicTargets() string {
	return "Targets"
}
