// Package cache package cache
package cache

import (
	"fmt"
	"trade_agent/pkg/sinopacapi"
)

// KeyOrderWaiting KeyOrderWaiting
func KeyOrderWaiting(stockNum string) string {
	return fmt.Sprintf("KeyOrderWaiting:%s", stockNum)
}

// KeyOrderBuy KeyOrderBuy
func KeyOrderBuy(stockNum string) string {
	return fmt.Sprintf("KeyOrderBuy:%s", stockNum)
}

// KeyOrderSell KeyOrderSell
func KeyOrderSell(stockNum string) string {
	return fmt.Sprintf("KeyOrderSell:%s", stockNum)
}

// KeyOrderSellFirst KeyOrderSellFirst
func KeyOrderSellFirst(stockNum string) string {
	return fmt.Sprintf("KeyOrderSellFirst:%s", stockNum)
}

// KeyOrderBuyLater KeyOrderBuyLater
func KeyOrderBuyLater(stockNum string) string {
	return fmt.Sprintf("KeyOrderBuyLater:%s", stockNum)
}

// GetOrderWaiting GetOrderWaiting
func (c *Cache) GetOrderWaiting(stockNum string) *sinopacapi.Order {
	defer c.lock.RUnlock()
	c.lock.RLock()
	if value, ok := c.Cache.Get(KeyOrderWaiting(stockNum)); ok {
		return value.(*sinopacapi.Order)
	}
	return nil
}

// GetOrderBuy GetOrderBuy
func (c *Cache) GetOrderBuy(stockNum string) []*sinopacapi.Order {
	defer c.lock.RUnlock()
	c.lock.RLock()
	if value, ok := c.Cache.Get(KeyOrderBuy(stockNum)); ok {
		return value.([]*sinopacapi.Order)
	}
	return nil
}

// GetOrderSell GetOrderSell
func (c *Cache) GetOrderSell(stockNum string) []*sinopacapi.Order {
	defer c.lock.RUnlock()
	c.lock.RLock()
	if value, ok := c.Cache.Get(KeyOrderSell(stockNum)); ok {
		return value.([]*sinopacapi.Order)
	}
	return nil
}

// GetOrderSellFirst GetOrderSellFirst
func (c *Cache) GetOrderSellFirst(stockNum string) []*sinopacapi.Order {
	defer c.lock.RUnlock()
	c.lock.RLock()
	if value, ok := c.Cache.Get(KeyOrderSellFirst(stockNum)); ok {
		return value.([]*sinopacapi.Order)
	}
	return nil
}

// GetOrderBuyLater GetOrderBuyLater
func (c *Cache) GetOrderBuyLater(stockNum string) []*sinopacapi.Order {
	defer c.lock.RUnlock()
	c.lock.RLock()
	if value, ok := c.Cache.Get(KeyOrderBuyLater(stockNum)); ok {
		return value.([]*sinopacapi.Order)
	}
	return nil
}
