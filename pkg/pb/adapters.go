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
		TickTime:  time.Unix(0, c.GetTs()),
		Close:     c.GetClose(),
		TickType:  c.GetTickType(),
		Volume:    c.GetVolume(),
		BidPrice:  c.GetBidPrice(),
		BidVolume: c.GetBidVolume(),
		AskPrice:  c.GetAskPrice(),
		AskVolume: c.GetAskVolume(),
		Open:      0,
		High:      0,
		Low:       0,
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
