// Package order package order
package order

import (
	"trade_agent/pkg/dbagent"
	"trade_agent/pkg/log"
)

func displayOrderResult(order *dbagent.OrderStatus) {
	var status string
	switch order.Status {
	case 1:
		status = "PendingSubmit"
	case 2:
		status = "PreSubmitted"
	case 3:
		status = "Submitted"
	case 4:
		status = "Failed"
	case 5:
		status = "Canceled"
	case 6:
		status = "Filled"
	case 7:
		status = "Filling"
	}
	log.Get().WithFields(map[string]interface{}{
		"Stock":    order.Stock.Number,
		"Action":   order.Action,
		"Quantity": order.Quantity,
		"Price":    order.Price,
	}).Infof("Order Status: %s", status)
}
