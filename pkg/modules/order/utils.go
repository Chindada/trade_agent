// Package order package order
package order

import (
	"trade_agent/pkg/dbagent"
	"trade_agent/pkg/log"
)

func displayOrderResult(order *dbagent.OrderStatus) {
	var status string
	switch order.Status {
	case 4:
		status = "Failed"
	case 5:
		status = "Canceled"
	case 6:
		status = "Filled"
	}
	log.Get().WithFields(map[string]interface{}{
		"Stock":    order.Stock.Number,
		"Action":   order.Action,
		"Quantity": order.Quantity,
		"Price":    order.Price,
	}).Errorf("Order Status: %s", status)
}
