// Package tickprocess package tickprocess
package tickprocess

import (
	"trade_agent/pkg/cache"
	"trade_agent/pkg/dbagent"
	"trade_agent/pkg/sinopacapi"
)

func realTimeTickArrAnalyzer(tick *dbagent.RealTimeTick, tickArr []*dbagent.RealTimeTick) sinopacapi.OrderAction {
	return 0
}

func realTimeBidAskArrAnalyzer(bidAsk *dbagent.RealTimeBidAsk, bidAskArr []*dbagent.RealTimeBidAsk) string {
	return ""
}

// HistoryTickAnalyzer HistoryTickAnalyzer
func HistoryTickAnalyzer(tickArr dbagent.HistoryTickArr) {
	var analyzeTickArr dbagent.HistoryTickArr
	for i, tick := range tickArr {
		if i == 0 {
			continue
		}
		if len(analyzeTickArr) > 1 {
			if analyzeTickArr.GetTotalTime() > 5 {
				var volumeSum int64
				for _, k := range analyzeTickArr {
					volumeSum += k.Volume
				}
				analyzeTickArr = []*dbagent.HistoryTick{}
				volumeArr := cache.GetCache().GetStockHistoryTickAnalyze(tick.Stock.Number)
				volumeArr = append(volumeArr, volumeSum)
				cache.GetCache().Set(cache.KeyStockHistoryTickAnalyze(tick.Stock.Number), volumeArr)
			}
		}
		analyzeTickArr = append(analyzeTickArr, tick)
	}
}
