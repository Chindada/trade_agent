// Package cache package cache
package cache

import (
	"fmt"
	"trade_agent/pkg/sinopacapi"
)

// KeyOrderWaiting KeyOrderWaiting
func KeyOrderWaiting(stockNum string) *Key {
	return &Key{
		Name: fmt.Sprintf("KeyOrderWaiting:%s", stockNum),
		Type: order,
	}
}

// GetOrderWaiting GetOrderWaiting
func (c *Cache) GetOrderWaiting(stockNum string) *sinopacapi.Order {
	c.lock.RLock()
	k := KeyOrderWaiting(stockNum)
	tmp := c.CacheMap[string(k.Type)]
	c.lock.RUnlock()
	if tmp == nil {
		return nil
	}
	if value, ok := tmp.Get(k.Name); ok {
		return value.(*sinopacapi.Order)
	}
	return nil
}

// KeyOrderBuy KeyOrderBuy
func KeyOrderBuy(stockNum string) *Key {
	return &Key{
		Name: fmt.Sprintf("KeyOrderBuy:%s", stockNum),
		Type: order,
	}
}

// GetOrderBuy GetOrderBuy
func (c *Cache) GetOrderBuy(stockNum string) []*sinopacapi.Order {
	c.lock.RLock()
	k := KeyOrderBuy(stockNum)
	tmp := c.CacheMap[string(k.Type)]
	c.lock.RUnlock()
	if tmp == nil {
		return []*sinopacapi.Order{}
	}
	if value, ok := tmp.Get(k.Name); ok {
		return value.([]*sinopacapi.Order)
	}
	return []*sinopacapi.Order{}
}

// AppendOrderBuy AppendOrderBuy
func (c *Cache) AppendOrderBuy(order *sinopacapi.Order) {
	tmp := c.GetOrderBuy(order.StockNum)
	tmp = append(tmp, order)
	c.Set(KeyOrderBuy(order.StockNum), tmp)
}

// KeyOrderSell KeyOrderSell
func KeyOrderSell(stockNum string) *Key {
	return &Key{
		Name: fmt.Sprintf("KeyOrderSell:%s", stockNum),
		Type: order,
	}
}

// GetOrderSell GetOrderSell
func (c *Cache) GetOrderSell(stockNum string) []*sinopacapi.Order {
	c.lock.RLock()
	k := KeyOrderSell(stockNum)
	tmp := c.CacheMap[string(k.Type)]
	c.lock.RUnlock()
	if tmp == nil {
		return []*sinopacapi.Order{}
	}
	if value, ok := tmp.Get(k.Name); ok {
		return value.([]*sinopacapi.Order)
	}
	return []*sinopacapi.Order{}
}

// AppendOrderSell AppendOrderSell
func (c *Cache) AppendOrderSell(order *sinopacapi.Order) {
	tmp := c.GetOrderSell(order.StockNum)
	tmp = append(tmp, order)
	c.Set(KeyOrderSell(order.StockNum), tmp)
}

// KeyOrderSellFirst KeyOrderSellFirst
func KeyOrderSellFirst(stockNum string) *Key {
	return &Key{
		Name: fmt.Sprintf("KeyOrderSellFirst:%s", stockNum),
		Type: order,
	}
}

// GetOrderSellFirst GetOrderSellFirst
func (c *Cache) GetOrderSellFirst(stockNum string) []*sinopacapi.Order {
	c.lock.RLock()
	k := KeyOrderSellFirst(stockNum)
	tmp := c.CacheMap[string(k.Type)]
	c.lock.RUnlock()
	if tmp == nil {
		return []*sinopacapi.Order{}
	}
	if value, ok := tmp.Get(k.Name); ok {
		return value.([]*sinopacapi.Order)
	}
	return []*sinopacapi.Order{}
}

// AppendOrderSellFirst AppendOrderSellFirst
func (c *Cache) AppendOrderSellFirst(order *sinopacapi.Order) {
	tmp := c.GetOrderSellFirst(order.StockNum)
	tmp = append(tmp, order)
	c.Set(KeyOrderSellFirst(order.StockNum), tmp)
}

// KeyOrderBuyLater KeyOrderBuyLater
func KeyOrderBuyLater(stockNum string) *Key {
	return &Key{
		Name: fmt.Sprintf("KeyOrderBuyLater:%s", stockNum),
		Type: order,
	}
}

// GetOrderBuyLater GetOrderBuyLater
func (c *Cache) GetOrderBuyLater(stockNum string) []*sinopacapi.Order {
	c.lock.RLock()
	k := KeyOrderBuyLater(stockNum)
	tmp := c.CacheMap[string(k.Type)]
	c.lock.RUnlock()
	if tmp == nil {
		return []*sinopacapi.Order{}
	}
	if value, ok := tmp.Get(k.Name); ok {
		return value.([]*sinopacapi.Order)
	}
	return []*sinopacapi.Order{}
}

// AppendOrderBuyLater AppendOrderBuyLater
func (c *Cache) AppendOrderBuyLater(order *sinopacapi.Order) {
	tmp := c.GetOrderBuyLater(order.StockNum)
	tmp = append(tmp, order)
	c.Set(KeyOrderBuyLater(order.StockNum), tmp)
}

// KeyOrderForward KeyOrderForward
func KeyOrderForward() *Key {
	return &Key{
		Name: "KeyOrderForward",
		Type: order,
	}
}

// GetOrderForward GetOrderForward
func (c *Cache) GetOrderForward() []*sinopacapi.Order {
	c.lock.RLock()
	k := KeyOrderForward()
	tmp := c.CacheMap[string(k.Type)]
	c.lock.RUnlock()
	if tmp == nil {
		return []*sinopacapi.Order{}
	}
	if value, ok := tmp.Get(k.Name); ok {
		return value.([]*sinopacapi.Order)
	}
	return []*sinopacapi.Order{}
}

// GetOrderForwardCountDetail GetOrderForwardCountDetail return remaining unfilled and total
func (c *Cache) GetOrderForwardCountDetail() (int64, int64) {
	var buy, sell int64
	for _, v := range c.GetOrderForward() {
		if v.Action == sinopacapi.ActionBuy {
			buy++
		} else {
			sell--
		}
	}
	return buy - sell, buy
}

// AppendOrderForward AppendOrderForward
func (c *Cache) AppendOrderForward(order *sinopacapi.Order) {
	tmp := c.GetOrderForward()
	tmp = append(tmp, order)
	c.Set(KeyOrderForward(), tmp)
}

// KeyOrderReverse KeyOrderReverse
func KeyOrderReverse() *Key {
	return &Key{
		Name: "KeyOrderReverse",
		Type: order,
	}
}

// GetOrderReverse GetOrderReverse
func (c *Cache) GetOrderReverse() []*sinopacapi.Order {
	c.lock.RLock()
	k := KeyOrderReverse()
	tmp := c.CacheMap[string(k.Type)]
	c.lock.RUnlock()
	if tmp == nil {
		return []*sinopacapi.Order{}
	}
	if value, ok := tmp.Get(k.Name); ok {
		return value.([]*sinopacapi.Order)
	}
	return []*sinopacapi.Order{}
}

// GetOrderReverseCountDetail GetOrderReverseCountDetail return remaining unfilled and total
func (c *Cache) GetOrderReverseCountDetail() (int64, int64) {
	var sellFirst, buyLater int64
	for _, v := range c.GetOrderReverse() {
		if v.Action == sinopacapi.ActionSellFirst {
			sellFirst++
		} else {
			buyLater--
		}
	}
	return sellFirst - buyLater, sellFirst
}

// AppendOrderReverse AppendOrderReverse
func (c *Cache) AppendOrderReverse(order *sinopacapi.Order) {
	tmp := c.GetOrderReverse()
	tmp = append(tmp, order)
	c.Set(KeyOrderReverse(), tmp)
}
