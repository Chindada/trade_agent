// Package tickprocess package tickprocess
package tickprocess

import (
	"trade_agent/pkg/dbagent"
	"trade_agent/pkg/sinopacapi"
)

func realTimeTickArrAnalyzer(tick *dbagent.RealTimeTick, tickArr []*dbagent.RealTimeTick) sinopacapi.OrderAction {
	return 0
}

func realTimeBidAskArrAnalyzer(bidAsk *dbagent.RealTimeBidAsk, bidAskArr []*dbagent.RealTimeBidAsk) string {
	return ""
}
