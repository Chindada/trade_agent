// Package order package order
package order

import (
	"trade_agent/pkg/eventbus"
	"trade_agent/pkg/log"
)

// InitOrder InitOrder
func InitOrder() {
	err := updateOrderStatus()
	if err != nil {
		log.Get().Panic(err)
	}
	err = eventbus.Get().Sub(eventbus.TopicStockOrderBuy(), stockOrderBuyCallback)
	if err != nil {
		log.Get().Panic(err)
	}
	err = eventbus.Get().Sub(eventbus.TopicStockOrderSellFirst(), stockOrderSellFirstCallback)
	if err != nil {
		log.Get().Panic(err)
	}
	err = eventbus.Get().Sub(eventbus.TopicStockOrderBuyLater(), stockOrderBuyLaterCallback)
	if err != nil {
		log.Get().Panic(err)
	}
	log.Get().Info("Initial Order")
}

func stockOrderBuyCallback() error {
	return nil
}

func stockOrderSellFirstCallback() error {
	return nil
}

func stockOrderBuyLaterCallback() error {
	return nil
}
