// Package tickprocess package tickprocess
package tickprocess

import (
	"trade_agent/pkg/config"
	"trade_agent/pkg/dbagent"
	"trade_agent/pkg/sinopacapi"
)

func realTimeTickArrAnalyzer(lastClose float64, tickArr dbagent.RealTimeTickArr, conf config.Analyze) sinopacapi.OrderAction {
	// stockNum := tickArr.GetStockNum()
	// historyTickAnalyze := cache.GetCache().GetStockHistoryTickAnalyze(stockNum)
	// bidAskStatus := cache.GetCache().GetRealTimeBidAskStatus(stockNum)
	// analyzeArr := tickArr.Analyzer()
	return 0
}

func realTimeBidAskArrAnalyzer(bidAsk *dbagent.RealTimeBidAsk, bidAskArr []*dbagent.RealTimeBidAsk) string {
	return ""
}
