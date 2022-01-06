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
		Type: historyTick,
	}
}

// GetStockHistoryTickAnalyze GetStockHistoryTickAnalyze
func (c *Cache) GetStockHistoryTickAnalyze(stockNum string) dbagent.AnalyzeVolumeArr {
	c.lock.RLock()
	k := KeyStockHistoryTickAnalyze(stockNum)
	tmp := c.CacheMap[string(k.Type)]
	c.lock.RUnlock()
	if tmp == nil {
		return dbagent.AnalyzeVolumeArr{}
	}
	if value, ok := tmp.Get(k.Name); ok {
		return value.(dbagent.AnalyzeVolumeArr)
	}
	return dbagent.AnalyzeVolumeArr{}
}

// AppendHistoryTickAnalyze AppendHistoryTickAnalyze
func (c *Cache) AppendHistoryTickAnalyze(stockNum string, arr dbagent.AnalyzeVolumeArr) {
	tmp := c.GetStockHistoryTickAnalyze(stockNum)
	tmp = append(tmp, arr...)
	c.Set(KeyOrderForward(), tmp)
}
