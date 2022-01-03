// Package balance package balance
package balance

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sort"
	"trade_agent/pkg/dbagent"
	"trade_agent/pkg/log"
	"trade_agent/pkg/routers/handlers"

	"github.com/gin-gonic/gin"
)

// AddHandlers AddHandlers
func AddHandlers(group *gin.RouterGroup) {
	group.GET("/balance", GetAllBalance)
	group.POST("/balance", ImportBalance)
	group.DELETE("/balance", DeletaAllBalance)
}

// GetAllBalance GetAllBalance
// @Summary GetAllBalance
// @tags Balance
// @accept json
// @produce json
// @success 200 {object} []dbagent.Balance
// @failure 500 {object} handlers.ErrorResponse
// @Router /balance [get]
func GetAllBalance(c *gin.Context) {
	var res handlers.ErrorResponse
	allBalance, err := dbagent.Get().GetAllBalance()
	if err != nil {
		log.Get().Error(err)
		res.Response = err.Error()
		c.JSON(http.StatusInternalServerError, res)
		return
	}
	if len(allBalance) > 1 {
		sort.Slice(allBalance, func(i, j int) bool {
			return allBalance[i].TradeDay.Before(allBalance[j].TradeDay)
		})
	}
	c.JSON(http.StatusOK, allBalance)
}

// ImportBalance ImportBalance
// @Summary ImportBalance
// @tags Balance
// @accept json
// @produce json
// @param body body []dbagent.Balance{} true "Body"
// @success 200
// @failure 500 {object} handlers.ErrorResponse
// @Router /balance [post]
func ImportBalance(c *gin.Context) {
	var res handlers.ErrorResponse
	body := []dbagent.Balance{}
	if byteArr, err := ioutil.ReadAll(c.Request.Body); err != nil {
		log.Get().Error(err)
		res.Response = err.Error()
		c.JSON(http.StatusInternalServerError, res)
		return
	} else if err := json.Unmarshal(byteArr, &body); err != nil {
		log.Get().Error(err)
		res.Response = err.Error()
		c.JSON(http.StatusInternalServerError, res)
		return
	}
	if len(body) > 1 {
		sort.Slice(body, func(i, j int) bool {
			return body[i].TradeDay.Before(body[j].TradeDay)
		})
	}
	for _, v := range body {
		tmp := v
		err := dbagent.Get().InsertOrUpdateBalance(&tmp)
		if err != nil {
			log.Get().Error(err)
			res.Response = err.Error()
			c.JSON(http.StatusInternalServerError, res)
			return
		}
	}
	c.JSON(http.StatusOK, nil)
}

// DeletaAllBalance DeletaAllBalance
// @Summary DeletaAllBalance
// @tags Balance
// @accept json
// @produce json
// @success 200
// @failure 500 {object} handlers.ErrorResponse
// @Router /balance [delete]
func DeletaAllBalance(c *gin.Context) {
	var res handlers.ErrorResponse
	if err := dbagent.Get().DeleteAllBalance(); err != nil {
		log.Get().Error(err)
		res.Response = err.Error()
		c.JSON(http.StatusInternalServerError, res)
		return
	}
	c.JSON(http.StatusOK, nil)
}
