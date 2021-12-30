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

// GetOrderWaiting GetOrderWaiting
func (c *Cache) GetOrderWaiting(stockNum string) *sinopacapi.Order {
	defer c.lock.RUnlock()
	c.lock.RLock()
	if value, ok := c.Cache.Get(KeyOrderWaiting(stockNum)); ok {
		return value.(*sinopacapi.Order)
	}
	return nil
}

// KeyOrderBuy KeyOrderBuy
func KeyOrderBuy(stockNum string) string {
	return fmt.Sprintf("KeyOrderBuy:%s", stockNum)
}

// GetOrderBuy GetOrderBuy
func (c *Cache) GetOrderBuy(stockNum string) []*sinopacapi.Order {
	defer c.lock.RUnlock()
	c.lock.RLock()
	if value, ok := c.Cache.Get(KeyOrderBuy(stockNum)); ok {
		return value.([]*sinopacapi.Order)
	}
	return []*sinopacapi.Order{}
}

// AppendOrderBuy AppendOrderBuy
func (c *Cache) AppendOrderBuy(order *sinopacapi.Order) {
	c.lock.RLock()
	var tmp []*sinopacapi.Order
	if value, ok := c.Cache.Get(KeyOrderBuy(order.StockNum)); ok {
		tmp = value.([]*sinopacapi.Order)
	}
	c.lock.RUnlock()
	tmp = append(tmp, order)
	c.Set(KeyOrderBuy(order.StockNum), tmp)
}

// KeyOrderSell KeyOrderSell
func KeyOrderSell(stockNum string) string {
	return fmt.Sprintf("KeyOrderSell:%s", stockNum)
}

// GetOrderSell GetOrderSell
func (c *Cache) GetOrderSell(stockNum string) []*sinopacapi.Order {
	defer c.lock.RUnlock()
	c.lock.RLock()
	if value, ok := c.Cache.Get(KeyOrderSell(stockNum)); ok {
		return value.([]*sinopacapi.Order)
	}
	return []*sinopacapi.Order{}
}

// AppendOrderSell AppendOrderSell
func (c *Cache) AppendOrderSell(order *sinopacapi.Order) {
	c.lock.RLock()
	var tmp []*sinopacapi.Order
	if value, ok := c.Cache.Get(KeyOrderSell(order.StockNum)); ok {
		tmp = value.([]*sinopacapi.Order)
	}
	c.lock.RUnlock()
	tmp = append(tmp, order)
	c.Set(KeyOrderSell(order.StockNum), tmp)
}

// KeyOrderSellFirst KeyOrderSellFirst
func KeyOrderSellFirst(stockNum string) string {
	return fmt.Sprintf("KeyOrderSellFirst:%s", stockNum)
}

// GetOrderSellFirst GetOrderSellFirst
func (c *Cache) GetOrderSellFirst(stockNum string) []*sinopacapi.Order {
	defer c.lock.RUnlock()
	c.lock.RLock()
	if value, ok := c.Cache.Get(KeyOrderSellFirst(stockNum)); ok {
		return value.([]*sinopacapi.Order)
	}
	return []*sinopacapi.Order{}
}

// AppendOrderSellFirst AppendOrderSellFirst
func (c *Cache) AppendOrderSellFirst(order *sinopacapi.Order) {
	c.lock.RLock()
	var tmp []*sinopacapi.Order
	if value, ok := c.Cache.Get(KeyOrderSellFirst(order.StockNum)); ok {
		tmp = value.([]*sinopacapi.Order)
	}
	c.lock.RUnlock()
	tmp = append(tmp, order)
	c.Set(KeyOrderSellFirst(order.StockNum), tmp)
}

// KeyOrderBuyLater KeyOrderBuyLater
func KeyOrderBuyLater(stockNum string) string {
	return fmt.Sprintf("KeyOrderBuyLater:%s", stockNum)
}

// GetOrderBuyLater GetOrderBuyLater
func (c *Cache) GetOrderBuyLater(stockNum string) []*sinopacapi.Order {
	defer c.lock.RUnlock()
	c.lock.RLock()
	if value, ok := c.Cache.Get(KeyOrderBuyLater(stockNum)); ok {
		return value.([]*sinopacapi.Order)
	}
	return []*sinopacapi.Order{}
}

// AppendOrderBuyLater AppendOrderBuyLater
func (c *Cache) AppendOrderBuyLater(order *sinopacapi.Order) {
	c.lock.RLock()
	var tmp []*sinopacapi.Order
	if value, ok := c.Cache.Get(KeyOrderBuyLater(order.StockNum)); ok {
		tmp = value.([]*sinopacapi.Order)
	}
	c.lock.RUnlock()
	tmp = append(tmp, order)
	c.Set(KeyOrderBuyLater(order.StockNum), tmp)
}

// KeyOrderForward KeyOrderForward
func KeyOrderForward() string {
	return "KeyOrderForward"
}

// GetOrderForward GetOrderForward
func (c *Cache) GetOrderForward() []*sinopacapi.Order {
	defer c.lock.RUnlock()
	c.lock.RLock()
	if value, ok := c.Cache.Get(KeyOrderForward()); ok {
		return value.([]*sinopacapi.Order)
	}
	return []*sinopacapi.Order{}
}

// GetOrderForwardCountDetail GetOrderForwardCountDetail return remaining unfilled and total
func (c *Cache) GetOrderForwardCountDetail() (int64, int64) {
	defer c.lock.RUnlock()
	c.lock.RLock()
	var buy, sell int64
	if value, ok := c.Cache.Get(KeyOrderForward()); ok {
		for _, v := range value.([]*sinopacapi.Order) {
			if v.Action == sinopacapi.ActionBuy {
				buy++
			} else {
				sell--
			}
		}
	}
	return buy - sell, buy
}

// AppendOrderForward AppendOrderForward
func (c *Cache) AppendOrderForward(order *sinopacapi.Order) {
	c.lock.RLock()
	var tmp []*sinopacapi.Order
	if value, ok := c.Cache.Get(KeyOrderForward()); ok {
		tmp = value.([]*sinopacapi.Order)
	}
	c.lock.RUnlock()
	tmp = append(tmp, order)
	c.Set(KeyOrderForward(), tmp)
}

// KeyOrderReverse KeyOrderReverse
func KeyOrderReverse() string {
	return "KeyOrderReverse"
}

// GetOrderReverse GetOrderReverse
func (c *Cache) GetOrderReverse() []*sinopacapi.Order {
	defer c.lock.RUnlock()
	c.lock.RLock()
	if value, ok := c.Cache.Get(KeyOrderReverse()); ok {
		return value.([]*sinopacapi.Order)
	}
	return []*sinopacapi.Order{}
}

// GetOrderReverseCountDetail GetOrderReverseCountDetail return remaining unfilled and total
func (c *Cache) GetOrderReverseCountDetail() (int64, int64) {
	defer c.lock.RUnlock()
	c.lock.RLock()
	var sellFirst, buyLater int64
	if value, ok := c.Cache.Get(KeyOrderReverse()); ok {
		for _, v := range value.([]*sinopacapi.Order) {
			if v.Action == sinopacapi.ActionSellFirst {
				sellFirst++
			} else {
				buyLater--
			}
		}
	}
	return sellFirst - buyLater, sellFirst
}

// AppendOrderReverse AppendOrderReverse
func (c *Cache) AppendOrderReverse(order *sinopacapi.Order) {
	c.lock.RLock()
	var tmp []*sinopacapi.Order
	if value, ok := c.Cache.Get(KeyOrderReverse()); ok {
		tmp = value.([]*sinopacapi.Order)
	}
	c.lock.RUnlock()
	tmp = append(tmp, order)
	c.Set(KeyOrderReverse(), tmp)
}
