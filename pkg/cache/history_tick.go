// Package cache package cache
package cache

import (
	"fmt"
	"trade_agent/pkg/dbagent"
)

// KeyStockHistoryTickAnalyze KeyStockHistoryTickAnalyze
func KeyStockHistoryTickAnalyze(stockNum string) *Key {
	return &Key{
		Name: fmt.Sprintf("KeyStockHistoryTickAnalyze:%s", stockNum),
		Type: keyTypeHistoryTickAnalyze(stockNum),
	}
}

// SetStockHistoryTickAnalyze SetStockHistoryTickAnalyze
func (c *Cache) SetStockHistoryTickAnalyze(stockNum string, arr dbagent.AnalyzeVolumeArr) {
	key := KeyStockHistoryTickAnalyze(stockNum)
	c.getCacheByType(key.Type).Set(key.Name, arr, noExpired)
}

// GetStockHistoryTickAnalyze GetStockHistoryTickAnalyze
func (c *Cache) GetStockHistoryTickAnalyze(stockNum string) dbagent.AnalyzeVolumeArr {
	k := KeyStockHistoryTickAnalyze(stockNum)
	if value, ok := c.getCacheByType(k.Type).Get(k.Name); ok {
		return value.(dbagent.AnalyzeVolumeArr)
	}
	return dbagent.AnalyzeVolumeArr{}
}

// AppendHistoryTickAnalyze AppendHistoryTickAnalyze
func (c *Cache) AppendHistoryTickAnalyze(stockNum string, arr dbagent.AnalyzeVolumeArr) {
	tmp := c.GetStockHistoryTickAnalyze(stockNum)
	tmp = append(tmp, arr...)
	c.SetStockHistoryTickAnalyze(stockNum, tmp)
}
