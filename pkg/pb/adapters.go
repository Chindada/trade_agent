// Package pb package pb
package pb

import (
	"time"
	"trade_agent/global"
	"trade_agent/pkg/cache"
	"trade_agent/pkg/dbagent"
	"trade_agent/pkg/log"
)

// ToStock ToStock
func (c *StockDetailMessage) ToStock() *dbagent.Stock {
	var dayTrade bool
	if c.GetDayTrade() == "Yes" {
		dayTrade = true
	}
	return &dbagent.Stock{
		Number:             c.GetCode(),
		Name:               c.GetName(),
		Exchange:           c.GetExchange(),
		Category:           c.GetCategory(),
		DayTrade:           dayTrade,
		LastClose:          c.GetReference(),
		LastVolume:         0,
		LastCloseChangePct: 0,
	}
}

// ToHistoryTick ToHistoryTick
func (c *HistoryTickMessage) ToHistoryTick(stockNum string) *dbagent.HistoryTick {
	return &dbagent.HistoryTick{
		Stock:     cache.GetCache().GetStock(stockNum),
		TickTime:  time.Unix(0, c.GetTs()).Add(-8 * time.Hour),
		Close:     c.GetClose(),
		TickType:  c.GetTickType(),
		Volume:    c.GetVolume(),
		BidPrice:  c.GetBidPrice(),
		BidVolume: c.GetBidVolume(),
		AskPrice:  c.GetAskPrice(),
		AskVolume: c.GetAskVolume(),
	}
}

// ToHistoryKbar ToHistoryKbar
func (c *HistoryKbarMessage) ToHistoryKbar(stockNum string) *dbagent.HistoryKbar {
	return &dbagent.HistoryKbar{
		Stock:    cache.GetCache().GetStock(stockNum),
		TickTime: time.Unix(0, c.GetTs()).Add(-8 * time.Hour),
		Close:    c.GetClose(),
		Open:     c.GetOpen(),
		High:     c.GetHigh(),
		Low:      c.GetLow(),
		Volume:   c.GetVolume(),
	}
}

// ToRealTimeTick ToRealTimeTick
func (c *RealTimeTickMessage) ToRealTimeTick() *dbagent.RealTimeTick {
	dataTime, err := time.ParseInLocation(global.LongTimeLayout, c.GetDateTime(), time.Local)
	if err != nil {
		return nil
	}
	return &dbagent.RealTimeTick{
		Stock:           cache.GetCache().GetStock(c.GetCode()),
		TickTime:        dataTime,
		Open:            c.GetOpen(),
		AvgPrice:        c.GetAvgPrice(),
		Close:           c.GetClose(),
		High:            c.GetHigh(),
		Low:             c.GetLow(),
		Amount:          c.GetAmount(),
		AmountSum:       c.GetTotalAmount(),
		Volume:          c.GetVolume(),
		VolumeSum:       c.GetTotalVolume(),
		TickType:        c.GetTickType(),
		ChgType:         c.GetChgType(),
		PriceChg:        c.GetPriceChg(),
		PctChg:          c.GetPctChg(),
		BidSideTotalVol: c.GetBidSideTotalVol(),
		AskSideTotalVol: c.GetAskSideTotalVol(),
		BidSideTotalCnt: c.GetBidSideTotalCnt(),
		AskSideTotalCnt: c.GetAskSideTotalCnt(),
		Suspend:         c.GetSuspend(),
		Simtrade:        c.GetSimtrade(),
	}
}

// ToTradeEvent ToTradeEvent
func (c *EventResponse) ToTradeEvent() *dbagent.CloudEvent {
	return &dbagent.CloudEvent{
		Event:     c.GetEvent(),
		EventCode: c.GetEventCode(),
		Info:      c.GetInfo(),
		Response:  c.GetRespCode(),
	}
}

// ToOrderStatus ToOrderStatus
func (c *OrderStatusHistoryMessage) ToOrderStatus() *dbagent.OrderStatus {
	actionMap := dbagent.ActionListMap
	statusMap := dbagent.StatusListMap
	orderTime, err := time.ParseInLocation(global.LongTimeLayout, c.GetOrderTime(), time.Local)
	if err != nil {
		log.Get().Error(err)
		return nil
	}
	return &dbagent.OrderStatus{
		Stock:     cache.GetCache().GetStock(c.GetCode()),
		Action:    actionMap[c.GetAction()],
		Price:     c.GetPrice(),
		Quantity:  c.GetQuantity(),
		Status:    statusMap[c.GetStatus()],
		OrderID:   c.GetOrderId(),
		OrderTime: orderTime,
	}
}

// ToRealTimeBidAsk ToRealTimeBidAsk
func (c *RealTimeBidAskMessage) ToRealTimeBidAsk() *dbagent.RealTimeBidAsk {
	dataTime, err := time.ParseInLocation(global.LongTimeLayout, c.GetDateTime(), time.Local)
	if err != nil {
		return nil
	}
	return &dbagent.RealTimeBidAsk{
		Stock:    cache.GetCache().GetStock(c.GetCode()),
		TickTime: dataTime,

		BidPrice1:   c.GetBidPrice()[0],
		BidVolume1:  c.GetBidVolume()[0],
		DiffBidVol1: c.GetDiffBidVol()[0],
		BidPrice2:   c.GetBidPrice()[1],
		BidVolume2:  c.GetBidVolume()[1],
		DiffBidVol2: c.GetDiffBidVol()[1],
		BidPrice3:   c.GetBidPrice()[2],
		BidVolume3:  c.GetBidVolume()[2],
		DiffBidVol3: c.GetDiffBidVol()[2],
		BidPrice4:   c.GetBidPrice()[3],
		BidVolume4:  c.GetBidVolume()[3],
		DiffBidVol4: c.GetDiffBidVol()[3],
		BidPrice5:   c.GetBidPrice()[4],
		BidVolume5:  c.GetBidVolume()[4],
		DiffBidVol5: c.GetDiffBidVol()[4],
		AskPrice1:   c.GetAskPrice()[0],
		AskVolume1:  c.GetAskVolume()[0],
		DiffAskVol1: c.GetDiffAskVol()[0],
		AskPrice2:   c.GetAskPrice()[1],
		AskVolume2:  c.GetAskVolume()[1],
		DiffAskVol2: c.GetDiffAskVol()[1],
		AskPrice3:   c.GetAskPrice()[2],
		AskVolume3:  c.GetAskVolume()[2],
		DiffAskVol3: c.GetDiffAskVol()[2],
		AskPrice4:   c.GetAskPrice()[3],
		AskVolume4:  c.GetAskVolume()[3],
		DiffAskVol4: c.GetDiffAskVol()[3],
		AskPrice5:   c.GetAskPrice()[4],
		AskVolume5:  c.GetAskVolume()[4],
		DiffAskVol5: c.GetDiffAskVol()[4],
		Suspend:     c.GetSuspend(),
		Simtrade:    c.GetSimtrade(),
	}
}
