package routers

import (
	"net/http"
	"strconv"

	"trade_agent/pkg/log"
	"trade_agent/pkg/sinopacapi"

	"github.com/gin-gonic/gin"
)

// AddDayTradeHandlersV1 AddDayTradeHandlersV1
func AddDayTradeHandlersV1(group *gin.RouterGroup) {
	group.GET("/day-trade/forward", CalculateDayTradeBalance)
}

// CalculateDayTradeBalance CalculateDayTradeBalance
// @Summary CalculateDayTradeBalance V1
// @tags DayTrade
// @accept json
// @produce json
// @param buy_price header string true "buy_price"
// @param buy_quantity header string true "buy_quantity"
// @param sell_price header string true "sell_price"
// @param sell_quantity header string true "sell_quantity"
// @success 200
// @failure 500 {object} ErrorResponse
// @Router /v1/day-trade/forward [get]
func CalculateDayTradeBalance(c *gin.Context) {
	var res ErrorResponse
	buyPriceString := c.Request.Header.Get("buy_price")
	buyPrice, err := strconv.ParseFloat(buyPriceString, 64)
	if err != nil {
		log.Get().Error(err)
		res.Response = err.Error()
		c.JSON(http.StatusInternalServerError, res)
		return
	}
	buyQuantityString := c.Request.Header.Get("buy_quantity")
	buyQuantity, err := strconv.ParseInt(buyQuantityString, 10, 64)
	if err != nil {
		log.Get().Error(err)
		res.Response = err.Error()
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	sellPriceString := c.Request.Header.Get("sell_price")
	sellPrice, err := strconv.ParseFloat(sellPriceString, 64)
	if err != nil {
		log.Get().Error(err)
		res.Response = err.Error()
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	sellQuantityString := c.Request.Header.Get("sell_quantity")
	sellQuantity, err := strconv.ParseInt(sellQuantityString, 10, 64)
	if err != nil {
		log.Get().Error(err)
		res.Response = err.Error()
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	pay := sinopacapi.GetStockBuyCost(buyPrice, buyQuantity)
	payDiscount := sinopacapi.GetStockTradeFeeDiscount(buyPrice, buyQuantity)
	earning := sinopacapi.GetStockSellCost(sellPrice, sellQuantity)
	earningDiscount := sinopacapi.GetStockTradeFeeDiscount(sellPrice, sellQuantity)

	c.JSON(http.StatusOK, gin.H{
		"balance": -pay + payDiscount + earning + earningDiscount,
	})
}
