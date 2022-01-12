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
		Type: waitingOrder,
	}
}

// SetOrderWaiting SetOrderWaiting
func (c *Cache) SetOrderWaiting(stockNum string, order *sinopacapi.Order) {
	key := KeyOrderWaiting(stockNum)
	c.getCacheByType(key.Type).Set(key.Name, order, noExpired)
}

// GetOrderWaiting GetOrderWaiting
func (c *Cache) GetOrderWaiting(stockNum string) *sinopacapi.Order {
	k := KeyOrderWaiting(stockNum)
	if value, ok := c.getCacheByType(k.Type).Get(k.Name); ok {
		return value.(*sinopacapi.Order)
	}
	return nil
}

// KeyOrderBuy KeyOrderBuy
func KeyOrderBuy(stockNum string) *Key {
	return &Key{
		Name: fmt.Sprintf("KeyOrderBuy:%s", stockNum),
		Type: buyOrder,
	}
}

// SetOrderBuy SetOrderBuy
func (c *Cache) SetOrderBuy(stockNum string, orderArr []*sinopacapi.Order) {
	key := KeyOrderBuy(stockNum)
	c.getCacheByType(key.Type).Set(key.Name, orderArr, noExpired)
}

// GetOrderBuy GetOrderBuy
func (c *Cache) GetOrderBuy(stockNum string) []*sinopacapi.Order {
	k := KeyOrderBuy(stockNum)
	if value, ok := c.getCacheByType(k.Type).Get(k.Name); ok {
		return value.([]*sinopacapi.Order)
	}
	return []*sinopacapi.Order{}
}

// AppendOrderBuy AppendOrderBuy
func (c *Cache) AppendOrderBuy(order *sinopacapi.Order) {
	tmp := c.GetOrderBuy(order.StockNum)
	tmp = append(tmp, order)
	c.SetOrderBuy(order.StockNum, tmp)
}

// KeyOrderSell KeyOrderSell
func KeyOrderSell(stockNum string) *Key {
	return &Key{
		Name: fmt.Sprintf("KeyOrderSell:%s", stockNum),
		Type: sellOrder,
	}
}

// SetOrderSell SetOrderSell
func (c *Cache) SetOrderSell(stockNum string, orderArr []*sinopacapi.Order) {
	key := KeyOrderSell(stockNum)
	c.getCacheByType(key.Type).Set(key.Name, orderArr, noExpired)
}

// GetOrderSell GetOrderSell
func (c *Cache) GetOrderSell(stockNum string) []*sinopacapi.Order {
	k := KeyOrderSell(stockNum)
	if value, ok := c.getCacheByType(k.Type).Get(k.Name); ok {
		return value.([]*sinopacapi.Order)
	}
	return []*sinopacapi.Order{}
}

// AppendOrderSell AppendOrderSell
func (c *Cache) AppendOrderSell(order *sinopacapi.Order) {
	tmp := c.GetOrderSell(order.StockNum)
	tmp = append(tmp, order)
	c.SetOrderSell(order.StockNum, tmp)
}

// KeyOrderSellFirst KeyOrderSellFirst
func KeyOrderSellFirst(stockNum string) *Key {
	return &Key{
		Name: fmt.Sprintf("KeyOrderSellFirst:%s", stockNum),
		Type: sellFirstOrder,
	}
}

// SetOrderSellFirst SetOrderSellFirst
func (c *Cache) SetOrderSellFirst(stockNum string, orderArr []*sinopacapi.Order) {
	key := KeyOrderSellFirst(stockNum)
	c.getCacheByType(key.Type).Set(key.Name, orderArr, noExpired)
}

// GetOrderSellFirst GetOrderSellFirst
func (c *Cache) GetOrderSellFirst(stockNum string) []*sinopacapi.Order {
	k := KeyOrderSellFirst(stockNum)
	if value, ok := c.getCacheByType(k.Type).Get(k.Name); ok {
		return value.([]*sinopacapi.Order)
	}
	return []*sinopacapi.Order{}
}

// AppendOrderSellFirst AppendOrderSellFirst
func (c *Cache) AppendOrderSellFirst(order *sinopacapi.Order) {
	tmp := c.GetOrderSellFirst(order.StockNum)
	tmp = append(tmp, order)
	c.SetOrderSellFirst(order.StockNum, tmp)
}

// KeyOrderBuyLater KeyOrderBuyLater
func KeyOrderBuyLater(stockNum string) *Key {
	return &Key{
		Name: fmt.Sprintf("KeyOrderBuyLater:%s", stockNum),
		Type: buyLaterOrder,
	}
}

// SetOrderBuyLater SetOrderBuyLater
func (c *Cache) SetOrderBuyLater(stockNum string, orderArr []*sinopacapi.Order) {
	key := KeyOrderBuyLater(stockNum)
	c.getCacheByType(key.Type).Set(key.Name, orderArr, noExpired)
}

// GetOrderBuyLater GetOrderBuyLater
func (c *Cache) GetOrderBuyLater(stockNum string) []*sinopacapi.Order {
	k := KeyOrderBuyLater(stockNum)
	if value, ok := c.getCacheByType(k.Type).Get(k.Name); ok {
		return value.([]*sinopacapi.Order)
	}
	return []*sinopacapi.Order{}
}

// AppendOrderBuyLater AppendOrderBuyLater
func (c *Cache) AppendOrderBuyLater(order *sinopacapi.Order) {
	tmp := c.GetOrderBuyLater(order.StockNum)
	tmp = append(tmp, order)
	c.SetOrderBuyLater(order.StockNum, tmp)
}

// KeyOrderForward KeyOrderForward
func KeyOrderForward() *Key {
	return &Key{
		Name: "KeyOrderForward",
		Type: forwardOrder,
	}
}

// SetOrderForward SetOrderForward
func (c *Cache) SetOrderForward(orderArr []*sinopacapi.Order) {
	key := KeyOrderForward()
	c.getCacheByType(key.Type).Set(key.Name, orderArr, noExpired)
}

// GetOrderForward GetOrderForward
func (c *Cache) GetOrderForward() []*sinopacapi.Order {
	k := KeyOrderForward()
	if value, ok := c.getCacheByType(k.Type).Get(k.Name); ok {
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
	c.SetOrderForward(tmp)
}

// KeyOrderReverse KeyOrderReverse
func KeyOrderReverse() *Key {
	return &Key{
		Name: "KeyOrderReverse",
		Type: reverseOrder,
	}
}

// SetOrderReverse SetOrderReverse
func (c *Cache) SetOrderReverse(orderArr []*sinopacapi.Order) {
	key := KeyOrderReverse()
	c.getCacheByType(key.Type).Set(key.Name, orderArr, noExpired)
}

// GetOrderReverse GetOrderReverse
func (c *Cache) GetOrderReverse() []*sinopacapi.Order {
	k := KeyOrderReverse()
	if value, ok := c.getCacheByType(k.Type).Get(k.Name); ok {
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
	c.SetOrderReverse(tmp)
}
