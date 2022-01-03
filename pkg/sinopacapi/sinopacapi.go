// Package sinopacapi package sinopacapi
package sinopacapi

import (
	"errors"
	"net/http"
	"strconv"
	"time"
	"trade_agent/pkg/config"
	"trade_agent/pkg/log"
	"trade_agent/pkg/restfulclient"
	"trade_agent/pkg/utils"

	"github.com/go-resty/resty/v2"
)

var globalClient *TradeAgent

// TradeAgent TradeAgent
type TradeAgent struct {
	Client    *resty.Client
	urlPrefix string
	token     string
}

// Order Order
type Order struct {
	StockNum  string      `json:"stock_num,omitempty" yaml:"stock_num"`
	Price     float64     `json:"price,omitempty" yaml:"price"`
	Quantity  int64       `json:"quantity,omitempty" yaml:"quantity"`
	OrderID   string      `json:"order_id,omitempty" yaml:"order_id"`
	Action    OrderAction `json:"action,omitempty" yaml:"action"`
	TradeTime time.Time   `json:"trade_time,omitempty" yaml:"trade_time"`
}

// InitSinpacAPI InitSinpacAPI
func InitSinpacAPI() {
	log.Get().Info("Initial SinopacAPI")

	serverConf := config.GetServerConfig()
	for {
		if utils.CheckPortIsOpen(serverConf.SinopacSRVHost, serverConf.SinopacSRVPort) {
			break
		}
		time.Sleep(time.Second)
	}
	mqConf := config.GetMQConfig()
	new := TradeAgent{
		Client:    restfulclient.Get(),
		urlPrefix: "http://" + serverConf.SinopacSRVHost + ":" + serverConf.SinopacSRVPort,
	}

	// check sinopac mq srv connect to mqtt broker
	err := new.AskSinpacMQSRVConnectMQ(mqConf)
	if err != nil {
		log.Get().Panic(err)
	}

	token, err := new.FetchServerToken()
	if err != nil {
		log.Get().Panic(err)
	} else {
		new.token = token
	}

	globalClient = &new
}

// Get Get
func Get() *TradeAgent {
	if globalClient == nil {
		log.Get().Panic("Trade Agent was not inititalized")
	}
	return globalClient
}

// GetToken GetToken
func (c *TradeAgent) GetToken() string {
	return c.token
}

// AskSinpacMQSRVConnectMQ AskSinpacMQSRVConnectMQ
func (c *TradeAgent) AskSinpacMQSRVConnectMQ(mqConf config.MQTT) (err error) {
	var resp *resty.Response
	resp, err = c.Client.R().
		SetBody(mqConf).
		SetResult(&ResponseCommon{}).
		Post(c.urlPrefix + urlAskSinpacMQSRVConnectMQ)
	if err != nil {
		return err
	} else if resp.StatusCode() != http.StatusOK {
		return errors.New("AskSinpacMQSRVConnectMQ API Fail")
	}
	if result := resp.Result().(*ResponseCommon).Result; result != StatusSuccuss {
		return errors.New(result)
	}
	return err
}

// FetchServerToken FetchServerToken
func (c *TradeAgent) FetchServerToken() (token string, err error) {
	var resp *resty.Response
	resp, err = c.Client.R().
		SetResult(&ResponseHealthStatus{}).
		Get(c.urlPrefix + urlFetchServerKey)
	if err != nil {
		return token, err
	} else if resp.StatusCode() != http.StatusOK {
		return token, errors.New("FetchServerKey API Fail")
	}
	if result := resp.Result().(*ResponseHealthStatus).Result; result != StatusSuccuss {
		return token, errors.New(result)
	}
	return resp.Result().(*ResponseHealthStatus).ServerToken, err
}

// RestartSinopacSRV RestartSinopacSRV
func (c *TradeAgent) RestartSinopacSRV() (err error) {
	var resp *resty.Response
	resp, err = c.Client.R().
		SetResult(&ResponseCommon{}).
		Get(c.urlPrefix + urlRestartSinopacSRV)
	if err != nil {
		return err
	} else if resp.StatusCode() != http.StatusOK {
		return errors.New("RestartSinopacSRV API Fail")
	}
	if result := resp.Result().(*ResponseCommon).Result; result != StatusSuccuss {
		return errors.New(result)
	}
	return err
}

// PlaceOrder PlaceOrder
func (c *TradeAgent) PlaceOrder(order Order) (res OrderResponse, err error) {
	var url string
	switch order.Action {
	case ActionBuy:
		url = urlPlaceOrderBuy
	case ActionSell:
		url = urlPlaceOrderSell
	case ActionSellFirst:
		url = urlPlaceOrderSellFirst
	case ActionBuyLater:
		url = urlPlaceOrderBuy
	}
	body := PlaceOrderBody{
		Stock:    order.StockNum,
		Price:    order.Price,
		Quantity: order.Quantity,
	}
	var resp *resty.Response
	resp, err = c.Client.R().
		SetBody(body).
		SetResult(&OrderResponse{}).
		Post(c.urlPrefix + url)
	if err != nil {
		return res, err
	} else if resp.StatusCode() != http.StatusOK {
		return res, errors.New("PlaceOrder API Fail")
	}
	return *resp.Result().(*OrderResponse), err
}

// CancelOrder CancelOrder
func (c *TradeAgent) CancelOrder(orderID string) (err error) {
	order := CancelOrderBody{
		OrderID: orderID,
	}
	var resp *resty.Response
	resp, err = c.Client.R().
		SetBody(order).
		SetResult(&ResponseCommon{}).
		Post(c.urlPrefix + urlCancelOrder)
	if err != nil {
		return err
	} else if resp.StatusCode() != http.StatusOK {
		return errors.New("CancelOrder API Fail")
	}
	switch resp.Result().(*ResponseCommon).Result {
	case StatusFail:
		return errors.New(StatusFail)
	case StatusAlreadyCanceled:
		return errors.New(StatusAlreadyCanceled)
	case StatusCancelOrderNotFound:
		return errors.New(StatusCancelOrderNotFound)
	}
	return err
}

// FetchOrderStatus FetchOrderStatus
func (c *TradeAgent) FetchOrderStatus() (err error) {
	var resp *resty.Response
	resp, err = c.Client.R().
		SetResult(&ResponseCommon{}).
		Get(c.urlPrefix + urlFetchOrderStatus)
	if err != nil {
		return err
	} else if resp.StatusCode() != http.StatusOK {
		return errors.New("FetchOrderStatus API Fail")
	}
	if result := resp.Result().(*ResponseCommon).Result; result != StatusSuccuss {
		return errors.New(result)
	}
	return err
}

// FetchHistoryCloseByStockArrDateArr FetchHistoryCloseByStockArrDateArr
func (c *TradeAgent) FetchHistoryCloseByStockArrDateArr(stockNumArr, dateArr []string) (err error) {
	stockAndDateArr := FetchHistoryCloseBody{
		StockNumArr: stockNumArr,
		DateArr:     dateArr,
	}
	var resp *resty.Response
	resp, err = c.Client.R().
		SetBody(stockAndDateArr).
		SetResult(&ResponseCommon{}).
		Post(c.urlPrefix + urlFetchHistoryCloseByStockArrDateArr)
	if err != nil {
		return err
	} else if resp.StatusCode() != http.StatusOK {
		return errors.New("FetchHistoryCloseByStockArrDateArr API Fail")
	}
	if result := resp.Result().(*ResponseCommon).Result; result != StatusSuccuss {
		return errors.New(result)
	}
	return err
}

// FetchAllSnapShot FetchAllSnapShot
func (c *TradeAgent) FetchAllSnapShot() (err error) {
	var resp *resty.Response
	resp, err = c.Client.R().
		SetResult(&ResponseCommon{}).
		Get(c.urlPrefix + urlFetchAllSnapShot)
	if err != nil {
		return err
	} else if resp.StatusCode() != http.StatusOK {
		return errors.New("FetchAllSnapShot API Fail")
	}
	if result := resp.Result().(*ResponseCommon).Result; result != StatusSuccuss {
		return errors.New(result)
	}
	return err
}

// FetchHistoryCloseByStockDateArr FetchHistoryCloseByStockDateArr
func (c *TradeAgent) FetchHistoryCloseByStockDateArr(stockNumArr []string, date string) (err error) {
	stockArr := FetchMultiDateHistoryCloseBody{
		StockNumArr: stockNumArr,
	}
	var resp *resty.Response
	resp, err = c.Client.R().
		SetHeader("X-Date", date).
		SetBody(stockArr).
		SetResult(&ResponseCommon{}).
		Post(c.urlPrefix + urlFetchHistoryCloseByStockDateArr)
	if err != nil {
		return err
	} else if resp.StatusCode() != http.StatusOK {
		return errors.New("FetchHistoryCloseByStockDateArr API Fail")
	}
	if result := resp.Result().(*ResponseCommon).Result; result != StatusSuccuss {
		return errors.New(result)
	}
	return err
}

// FetchHistoryTSECloseByDate FetchHistoryTSECloseByDate
func (c *TradeAgent) FetchHistoryTSECloseByDate(date string) (err error) {
	var resp *resty.Response
	resp, err = c.Client.R().
		SetHeader("X-Date", date).
		SetResult(&ResponseCommon{}).
		Post(c.urlPrefix + urlFetchHistoryTSECloseByDate)
	if err != nil {
		return err
	} else if resp.StatusCode() != http.StatusOK {
		return errors.New("FetchHistoryTSECloseByDate API Fail")
	}
	if result := resp.Result().(*ResponseCommon).Result; result != StatusSuccuss {
		return errors.New(result)
	}
	return err
}

// FetchVolumeRankByDate FetchVolumeRankByDate
func (c *TradeAgent) FetchVolumeRankByDate(date string, count int64) (err error) {
	var resp *resty.Response
	resp, err = c.Client.R().
		SetHeader("X-Count", strconv.FormatInt(count, 10)).
		SetHeader("X-Date", date).
		SetResult(&ResponseCommon{}).
		Get(c.urlPrefix + urlFetchVolumeRankByDate)
	if err != nil {
		return err
	} else if resp.StatusCode() != http.StatusOK {
		return errors.New("FetchVolumeRankByDate API Fail")
	}
	if result := resp.Result().(*ResponseCommon).Result; result != StatusSuccuss {
		return errors.New(result)
	}
	return err
}

// FetchHistoryKbarByDateRange FetchHistoryKbarByDateRange
func (c *TradeAgent) FetchHistoryKbarByDateRange(stockNum string, start, end string) (err error) {
	stockAndDateArr := FetchHistoryKbarBody{
		StockNum:  stockNum,
		StartDate: start,
		EndDate:   end,
	}
	resp, err := c.Client.R().
		SetBody(stockAndDateArr).
		SetResult(&ResponseCommon{}).
		Post(c.urlPrefix + urlFetchHistoryKbarByDateRange)
	if err != nil {
		return err
	} else if resp.StatusCode() != http.StatusOK {
		return errors.New("FetchHistoryKbarByDateRange API Fail")
	}
	if result := resp.Result().(*ResponseCommon).Result; result != StatusSuccuss {
		return errors.New(result)
	}
	return err
}

// FetchHistoryTSEKbarByDate FetchHistoryTSEKbarByDate
func (c *TradeAgent) FetchHistoryTSEKbarByDate(date string) (err error) {
	stockAndDateArr := FetchHistoryKbarBody{
		StartDate: date,
		EndDate:   date,
	}
	resp, err := c.Client.R().
		SetBody(stockAndDateArr).
		SetResult(&ResponseCommon{}).
		Post(c.urlPrefix + urlFetchHistoryTSEKbarByDate)
	if err != nil {
		return err
	} else if resp.StatusCode() != http.StatusOK {
		return errors.New("FetchHistoryTSEKbarByDate API Fail")
	}
	if result := resp.Result().(*ResponseCommon).Result; result != StatusSuccuss {
		return errors.New(result)
	}
	return err
}

// FetchHistoryTickByStockAndDate FetchHistoryTickByStockAndDate
func (c *TradeAgent) FetchHistoryTickByStockAndDate(stockNum, date string) (err error) {
	stockAndDate := FetchHistoryTickBody{
		StockNum: stockNum,
		Date:     date,
	}
	var resp *resty.Response
	resp, err = c.Client.R().
		SetBody(stockAndDate).
		SetResult(&ResponseCommon{}).
		Post(c.urlPrefix + urlFetchHistoryTickByStockAndDate)
	if err != nil {
		return err
	} else if resp.StatusCode() != http.StatusOK {
		return errors.New("FetchHistoryTickByStockAndDate API Fail")
	}
	if result := resp.Result().(*ResponseCommon).Result; result != StatusSuccuss {
		return errors.New(result)
	}
	return err
}

// FetchAllStockDetail FetchAllStockDetail
func (c *TradeAgent) FetchAllStockDetail() (err error) {
	var resp *resty.Response
	resp, err = c.Client.R().
		SetResult(&ResponseCommon{}).
		Get(c.urlPrefix + urlFetchAllStockDetail)
	if err != nil {
		return err
	} else if resp.StatusCode() != http.StatusOK {
		return errors.New("FetchAllStockDetail API Fail")
	}
	if result := resp.Result().(*ResponseCommon).Result; result != StatusSuccuss {
		return errors.New(result)
	}
	return err
}

// SubRealTimeTick SubRealTimeTick
func (c *TradeAgent) SubRealTimeTick(stockArr []string) (err error) {
	stocks := SubscribeBody{
		StockNumArr: stockArr,
	}
	var resp *resty.Response
	resp, err = c.Client.R().
		SetBody(stocks).
		SetResult(&ResponseCommon{}).
		Post(c.urlPrefix + urlSubRealTimeTick)
	if err != nil {
		return err
	} else if resp.StatusCode() != http.StatusOK {
		return errors.New("SubRealTimeTick API Fail")
	}
	if result := resp.Result().(*ResponseCommon).Result; result != StatusSuccuss {
		return errors.New(result)
	}
	return err
}

// UnSubRealTimeTick UnSubRealTimeTick
func (c *TradeAgent) UnSubRealTimeTick(stockArr []string) (err error) {
	stocks := SubscribeBody{
		StockNumArr: stockArr,
	}
	var resp *resty.Response
	resp, err = c.Client.R().
		SetBody(stocks).
		SetResult(&ResponseCommon{}).
		Post(c.urlPrefix + urlUnSubRealTimeTickByStock)
	if err != nil {
		return err
	} else if resp.StatusCode() != http.StatusOK {
		return errors.New("UnSubRealTimeTick API Fail")
	}
	if result := resp.Result().(*ResponseCommon).Result; result != StatusSuccuss {
		return errors.New(result)
	}
	return err
}

// SubBidAsk SubBidAsk
func (c *TradeAgent) SubBidAsk(stockArr []string) (err error) {
	stocks := SubscribeBody{
		StockNumArr: stockArr,
	}
	var resp *resty.Response
	resp, err = c.Client.R().
		SetBody(stocks).
		SetResult(&ResponseCommon{}).
		Post(c.urlPrefix + urlSubBidAsk)
	if err != nil {
		return err
	} else if resp.StatusCode() != http.StatusOK {
		return errors.New("SubBidAsk API Fail")
	}
	if result := resp.Result().(*ResponseCommon).Result; result != StatusSuccuss {
		return errors.New(result)
	}
	return err
}

// UnSubBidAsk UnSubBidAsk
func (c *TradeAgent) UnSubBidAsk(stockArr []string) (err error) {
	stocks := SubscribeBody{
		StockNumArr: stockArr,
	}
	var resp *resty.Response
	resp, err = c.Client.R().
		SetBody(stocks).
		SetResult(&ResponseCommon{}).
		Post(c.urlPrefix + urlUnSubBidAskByStock)
	if err != nil {
		return err
	} else if resp.StatusCode() != http.StatusOK {
		return errors.New("UnSubBidAsk API Fail")
	}
	if result := resp.Result().(*ResponseCommon).Result; result != StatusSuccuss {
		return errors.New(result)
	}
	return err
}

// UnSubscribeAllByType UnSubscribeAllByType
func (c *TradeAgent) UnSubscribeAllByType(dataType TickType) (err error) {
	var url string
	switch {
	case dataType == TickTypeStockRealTime:
		url = urlUnSubscribeAllRealTimeTick
	case dataType == TickTypeStockBidAsk:
		url = urlUnSubscribeAllBidAsk
	}
	var resp *resty.Response
	resp, err = c.Client.R().
		SetResult(&ResponseCommon{}).
		Get(c.urlPrefix + url)
	if err != nil {
		return err
	} else if resp.StatusCode() != http.StatusOK {
		return errors.New("UnSubscribeAllByType API Fail")
	}
	if result := resp.Result().(*ResponseCommon).Result; result != StatusSuccuss {
		return errors.New(result)
	}
	return err
}
