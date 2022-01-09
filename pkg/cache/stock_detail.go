// Package cache package cache
package cache

import (
	"fmt"
	"trade_agent/pkg/dbagent"
)

// KeyStockDetail KeyStockDetail
func KeyStockDetail(stockNum string) *Key {
	return &Key{
		Name: fmt.Sprintf("KeyStockDetail:%s", stockNum),
		Type: stockDetail,
	}
}

// SetStockDetail SetStockDetail
func (c *Cache) SetStockDetail(stock *dbagent.Stock) {
	key := KeyStockDetail(stock.Number)
	c.getCacheByType(key.Type).Set(key.Name, stock, noExpired)
}

// GetStock GetStock
func (c *Cache) GetStock(stockNum string) *dbagent.Stock {
	c.lock.RLock()
	k := KeyStockDetail(stockNum)
	tmp := c.CacheMap[string(k.Type)]
	c.lock.RUnlock()
	if tmp == nil {
		return nil
	}
	if value, ok := tmp.Get(k.Name); ok {
		return value.(*dbagent.Stock)
	}
	return nil
}

// GetStockID GetStockID
func (c *Cache) GetStockID(stockNum string) int64 {
	return int64(c.GetStock(stockNum).ID)
}
