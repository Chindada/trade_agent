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

	"github.com/go-resty/resty/v2"
)

var globalClient *TradeAgent

// TradeAgent TradeAgent
type TradeAgent struct {
	Client    *resty.Client
	urlPrefix string
}

// Order Order
type Order struct {
	StockNum  string
	Price     float64
	Quantity  int64
	OrderID   string
	Action    OrderAction
	TradeTime time.Time
}

// Get Get
func Get() *TradeAgent {
	if globalClient == nil {
		log.Get().Panic("Trade Agent was not inititalized")
	}
	return globalClient
}

// InitSinpacAPI InitSinpacAPI
func InitSinpacAPI() {
	log.Get().Info("Initial SinopacAPI")
	conf, err := config.Get()
	if err != nil {
		log.Get().Panic(err)
	}
	serverConf := conf.GetServerConfig()
	mqConf := conf.GetMQConfig()
	new := TradeAgent{
		Client:    restfulclient.Get(),
		urlPrefix: "http://" + serverConf.SinopacSRVHost + ":" + serverConf.SinopacSRVPort,
	}
	// check sinopac mq srv connect to mqtt broker
	err = new.AskSinpacMQSRVConnectMQ(mqConf)
	if err != nil {
		log.Get().Panic(err)
	}
	globalClient = &new
}

// AskSinpacMQSRVConnectMQ AskSinpacMQSRVConnectMQ
func (c *TradeAgent) AskSinpacMQSRVConnectMQ(mqConf config.MQTT) (err error) {
	var resp *resty.Response
	resp, err = c.Client.R().
		SetBody(mqConf).
		SetResult(&ResponseCommon{}).
		Post(c.urlPrefix + urlAskSinpacConnectMQ)
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

// FetchServerKey FetchServerKey
func (c *TradeAgent) FetchServerKey() (token string, err error) {
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
	}
	body := OrderBody{
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
	order := OrderCancelBody{
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

// FetchLastCloseByStockArrDateArr FetchLastCloseByStockArrDateArr
func (c *TradeAgent) FetchLastCloseByStockArrDateArr(stockNumArr, dateArr []string) (err error) {
	stockAndDateArr := FetchLastCloseBody{
		StockNumArr: stockNumArr,
		DateArr:     dateArr,
	}
	var resp *resty.Response
	resp, err = c.Client.R().
		SetBody(stockAndDateArr).
		SetResult(&ResponseCommon{}).
		Post(c.urlPrefix + urlFetchLastCloseByStockArrDateArr)
	if err != nil {
		return err
	} else if resp.StatusCode() != http.StatusOK {
		return errors.New("FetchLastCloseByStockArrDateArr API Fail")
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

// FetchStockCloseMapByStockDateArr FetchStockCloseMapByStockDateArr
func (c *TradeAgent) FetchStockCloseMapByStockDateArr(stockNumArr []string, dateArr []time.Time) (err error) {
	stockArr := FetchLastCountBody{
		StockNumArr: stockNumArr,
	}
	for _, date := range dateArr {
		var resp *resty.Response
		resp, err = c.Client.R().
			SetHeader("X-Date", date.Format(shortTimeLayout)).
			SetBody(stockArr).
			SetResult(&ResponseCommon{}).
			Post(c.urlPrefix + urlFetchStockCloseMapByStockDateArr)
		if err != nil {
			return err
		} else if resp.StatusCode() != http.StatusOK {
			return errors.New("FetchStockCloseMapByStockDateArr API Fail")
		}
		if result := resp.Result().(*ResponseCommon).Result; result != StatusSuccuss {
			return errors.New(result)
		}
	}
	return err
}

// FetchTSE001CloseByDate FetchTSE001CloseByDate
func (c *TradeAgent) FetchTSE001CloseByDate(date time.Time) (err error) {
	var resp *resty.Response
	resp, err = c.Client.R().
		SetHeader("X-Date", date.Format(shortTimeLayout)).
		SetResult(&ResponseCommon{}).
		Post(c.urlPrefix + urlFetchTSE001CloseByDate)
	if err != nil {
		return err
	} else if resp.StatusCode() != http.StatusOK {
		return errors.New("FetchTSE001CloseByDate API Fail")
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

// FetchKbarByDateRange FetchKbarByDateRange
func (c *TradeAgent) FetchKbarByDateRange(stockNum string, start, end time.Time) (err error) {
	stockAndDateArr := FetchKbarBody{
		StockNum:  stockNum,
		StartDate: start.Format(shortTimeLayout),
		EndDate:   end.Format(shortTimeLayout),
	}
	resp, err := c.Client.R().
		SetBody(stockAndDateArr).
		SetResult(&ResponseCommon{}).
		Post(c.urlPrefix + urlFetchKbarByDateRange)
	if err != nil {
		return err
	} else if resp.StatusCode() != http.StatusOK {
		return errors.New("FetchKbarByDateRange API Fail")
	}
	if result := resp.Result().(*ResponseCommon).Result; result != StatusSuccuss {
		return errors.New(result)
	}
	return err
}

// FetchTSE001KbarByDate FetchTSE001KbarByDate
func (c *TradeAgent) FetchTSE001KbarByDate(date time.Time) (err error) {
	stockAndDateArr := FetchKbarBody{
		StartDate: date.Format(shortTimeLayout),
		EndDate:   date.Format(shortTimeLayout),
	}
	resp, err := c.Client.R().
		SetBody(stockAndDateArr).
		SetResult(&ResponseCommon{}).
		Post(c.urlPrefix + urlFetchTSE001KbarByDate)
	if err != nil {
		return err
	} else if resp.StatusCode() != http.StatusOK {
		return errors.New("FetchTSEKbarByDateRange API Fail")
	}
	if result := resp.Result().(*ResponseCommon).Result; result != StatusSuccuss {
		return errors.New(result)
	}
	return err
}

// FetchEntireTickByStockAndDate FetchEntireTickByStockAndDate
func (c *TradeAgent) FetchEntireTickByStockAndDate(stockNum, date string) (err error) {
	stockAndDate := FetchBody{
		StockNum: stockNum,
		Date:     date,
	}
	var resp *resty.Response
	resp, err = c.Client.R().
		SetBody(stockAndDate).
		SetResult(&ResponseCommon{}).
		Post(c.urlPrefix + urlFetchEntireTickByStockAndDate)
	if err != nil {
		return err
	} else if resp.StatusCode() != http.StatusOK {
		return errors.New("FetchEntireTickByStockAndDate API Fail")
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

// SubStreamTick SubStreamTick
func (c *TradeAgent) SubStreamTick(stockArr []string) (err error) {
	stocks := SubscribeBody{
		StockNumArr: stockArr,
	}
	var resp *resty.Response
	resp, err = c.Client.R().
		SetBody(stocks).
		SetResult(&ResponseCommon{}).
		Post(c.urlPrefix + urlSubStreamTick)
	if err != nil {
		return err
	} else if resp.StatusCode() != http.StatusOK {
		return errors.New("SubStreamTick API Fail")
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

// UnSubscribeAllByType UnSubscribeAllByType
func (c *TradeAgent) UnSubscribeAllByType(dataType TickType) (err error) {
	var url string
	switch {
	case dataType == StreamType:
		url = urlUnSubscribeAllStream
	case dataType == BidAskType:
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
