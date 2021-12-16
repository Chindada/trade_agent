// Package eventbus package eventbus
package eventbus

import "fmt"

// TopicStockDetail TopicStockDetail
func TopicStockDetail(stockNum string) string {
	return fmt.Sprintf("StockDetail:%s", stockNum)
}

// TopicTargets TopicTargets
func TopicTargets() string {
	return "Targets"
}
