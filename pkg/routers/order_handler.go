// Package routers package routers
package routers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sort"

	"trade_agent/pkg/dbagent"
	"trade_agent/pkg/log"

	"github.com/gin-gonic/gin"
)

// AddOrderHandlersV1 AddOrderHandlersV1
func AddOrderHandlersV1(group *gin.RouterGroup) {
	group.GET("/order", GetAllOrder)
	group.POST("/order", ImportOrder)
	group.DELETE("/order", DeletaAllOrder)
}

// GetAllOrder GetAllOrder
// @Summary GetAllOrder V1
// @tags Order
// @accept json
// @produce json
// @success 200 {object} []dbagent.OrderStatus
// @failure 500 {object} ErrorResponse
// @Router /v1/order [get]
func GetAllOrder(c *gin.Context) {
	var res ErrorResponse
	allOrder, err := dbagent.Get().GetAllOrderStatus()
	if err != nil {
		log.Get().Error(err)
		res.Response = err.Error()
		c.JSON(http.StatusInternalServerError, res)
		return
	}
	if len(allOrder) > 1 {
		sort.Slice(allOrder, func(i, j int) bool {
			return allOrder[i].OrderTime.Before(allOrder[j].OrderTime)
		})
	}
	c.JSON(http.StatusOK, allOrder)
}

// ImportOrder ImportOrder
// @Summary ImportOrder V1
// @tags Order
// @accept json
// @produce json
// @param body body []dbagent.OrderStatus{} true "Body"
// @success 200
// @failure 500 {object} ErrorResponse
// @Router /v1/order [post]
func ImportOrder(c *gin.Context) {
	var res ErrorResponse
	body := []dbagent.OrderStatus{}
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
			return body[i].OrderTime.Before(body[j].OrderTime)
		})
	}
	var tmp []*dbagent.OrderStatus
	for _, v := range body {
		order := v
		tmp = append(tmp, &order)
	}
	err := dbagent.Get().InsertOrUpdateMultiOrderStatus(tmp)
	if err != nil {
		log.Get().Error(err)
		res.Response = err.Error()
		c.JSON(http.StatusInternalServerError, res)
		return
	}
	c.JSON(http.StatusOK, nil)
}

// DeletaAllOrder DeletaAllOrder
// @Summary DeletaAllOrder V1
// @tags Order
// @accept json
// @produce json
// @success 200
// @failure 500 {object} ErrorResponse
// @Router /v1/order [delete]
func DeletaAllOrder(c *gin.Context) {
	var res ErrorResponse
	if err := dbagent.Get().DeleteAllOrderStatus(); err != nil {
		log.Get().Error(err)
		res.Response = err.Error()
		c.JSON(http.StatusInternalServerError, res)
		return
	}
	c.JSON(http.StatusOK, nil)
}
