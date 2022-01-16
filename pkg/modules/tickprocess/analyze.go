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

func realTimeTickArrActionGenerator(tickArr, lastPeriodArr dbagent.RealTimeTickArr, conf config.Analyze) sinopacapi.OrderAction {
	lastTick := tickArr.GetLastTick()
	if lastTick == nil {
		return 0
	} else if lastTick.PctChg < conf.CloseChangeRatioLow || lastTick.PctChg > conf.CloseChangeRatioHigh {
		return 0
	}

	stockNum := tickArr.GetStockNum()
	if stockNum == "" {
		return 0
	}

	historyTickAnalyze := cache.GetCache().GetStockHistoryTickAnalyze(stockNum)
	if pr := historyTickAnalyze.GetPRByVolume(lastPeriodArr.GetTotalVolume()); pr < conf.VolumePRLow || pr > conf.VolumePRHigh {
		return 0
	}

	outInRatio := lastPeriodArr.GetOutInRatio()
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
	return getTradeOutAction(availableActionMap, conf, rsi)
}

func getTradeOutAction(availableActionMap map[sinopacapi.OrderAction]bool, conf config.Analyze, rsi float64) sinopacapi.OrderAction {
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
	forwardRest := len(historyOrderBuy) - len(historyOrderSell)

	historyOrderSellFirst := cache.GetCache().GetOrderSellFirst(stockNum)
	historyOrderBuyLater := cache.GetCache().GetOrderBuyLater(stockNum)
	reverseRest := len(historyOrderSellFirst) - len(historyOrderBuyLater)

	if forwardRest != 0 {
		tmp[sinopacapi.ActionSell] = true
		preTime = historyOrderBuy[len(historyOrderBuy)-forwardRest].TradeTime
	} else if reverseRest == 0 {
		tmp[sinopacapi.ActionBuy] = true
	}

	if reverseRest != 0 {
		tmp[sinopacapi.ActionBuyLater] = true
		preTime = historyOrderSellFirst[len(historyOrderSellFirst)-reverseRest].TradeTime
	} else if forwardRest == 0 {
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
