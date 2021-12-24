// Package order package order
package order

import (
	"trade_agent/pkg/eventbus"
	"trade_agent/pkg/log"
	"trade_agent/pkg/sinopacapi"
)

// InitOrder InitOrder
func InitOrder() {
	err := updateOrderStatus()
	if err != nil {
		log.Get().Panic(err)
	}

	err = eventbus.Get().Sub(eventbus.TopicStockOrder(), orderBuyCallback)
	if err != nil {
		log.Get().Panic(err)
	}
	log.Get().Info("Initial Order")
}

func orderBuyCallback(order sinopacapi.Order) error {
	res, err := sinopacapi.Get().PlaceOrder(order)
	if err != nil {
		return err
	}
	log.Get().Infof("Order ID %s", res.OrderID)
	return nil
}
