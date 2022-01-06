// Package tickprocess package tickprocess
package tickprocess

import (
	"time"
	"trade_agent/pkg/cache"
	"trade_agent/pkg/config"
	"trade_agent/pkg/dbagent"
	"trade_agent/pkg/sinopacapi"
)

func realTimeBidAskArrStatusGenerator(bidAsk *dbagent.RealTimeBidAsk, bidAskArr []*dbagent.RealTimeBidAsk) string {
	return ""
}

func realTimeTickArrActionGenerator(tickArr dbagent.RealTimeTickArr, conf config.Analyze) sinopacapi.OrderAction {
	lastTick := tickArr.GetLastTick()
	if lastTick == nil {
		return 0
	}
	if lastTick.PctChg < conf.CloseChangeRatioLow || lastTick.PctChg > conf.CloseChangeRatioHigh {
		return 0
	}

	stockNum := tickArr.GetStockNum()
	historyTickAnalyze := cache.GetCache().GetStockHistoryTickAnalyze(stockNum)
	pr := historyTickAnalyze.GetPRByVolume(tickArr.GetLastPeriodVolume())
	// if pr < conf.VolumePR {
	// 	return 0
	// }
	outInRatio := tickArr.GetOutInRatio()
	tmp := &dbagent.RealTimeTickAnalyze{
		Stock:      lastTick.Stock,
		TickTime:   lastTick.TickTime,
		PR:         pr,
		OutInRatio: outInRatio,
	}
	_ = dbagent.Get().InsertRealTimeTickAnalyze(tmp)

	availableActionMap, preTime := getAvailableAction(stockNum)
	if outInRatio > conf.OutInRatio && availableActionMap[sinopacapi.ActionBuy] {
		return sinopacapi.ActionBuy
	}
	if 100-outInRatio < conf.InOutRatio && availableActionMap[sinopacapi.ActionSellFirst] {
		return sinopacapi.ActionSellFirst
	}

	rsi := tickArr.GetRSIByTickTime(preTime, conf.RSIMinCount)
	if rsi == 0 {
		return 0
	}

	if availableActionMap[sinopacapi.ActionSell] && rsi > conf.RSIHigh {
		return sinopacapi.ActionSell
	}

	if availableActionMap[sinopacapi.ActionBuyLater] && rsi < conf.RSILow {
		return sinopacapi.ActionBuyLater
	}

	return 0
}

func getAvailableAction(stockNum string) (map[sinopacapi.OrderAction]bool, time.Time) {
	preTime := time.Time{}
	tmp := make(map[sinopacapi.OrderAction]bool)

	historyOrderBuy := cache.GetCache().GetOrderBuy(stockNum)
	historyOrderSell := cache.GetCache().GetOrderSell(stockNum)
	if len(historyOrderBuy) > len(historyOrderSell) {
		tmp[sinopacapi.ActionSell] = true
		preTime = historyOrderBuy[len(historyOrderBuy)-1].TradeTime
	} else {
		tmp[sinopacapi.ActionBuy] = true
	}

	historyOrderSellFirst := cache.GetCache().GetOrderSellFirst(stockNum)
	historyOrderBuyLater := cache.GetCache().GetOrderBuyLater(stockNum)
	if len(historyOrderSellFirst) > len(historyOrderBuyLater) {
		tmp[sinopacapi.ActionBuyLater] = true
		preTime = historyOrderSellFirst[len(historyOrderSellFirst)-1].TradeTime
	} else {
		tmp[sinopacapi.ActionSellFirst] = true
	}

	// check priority
	if tmp[sinopacapi.ActionBuy] && tmp[sinopacapi.ActionBuyLater] {
		delete(tmp, sinopacapi.ActionBuy)
	}
	if tmp[sinopacapi.ActionSell] && tmp[sinopacapi.ActionSellFirst] {
		delete(tmp, sinopacapi.ActionSellFirst)
	}
	return tmp, preTime
}
