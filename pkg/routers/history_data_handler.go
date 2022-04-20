// Package routers package routers
package routers

import (
	"net/http"
	"strconv"
	"time"

	"trade_agent/global"
	"trade_agent/pkg/cache"
	"trade_agent/pkg/dbagent"
	"trade_agent/pkg/eventbus"

	"github.com/gin-gonic/gin"
)

// AddHistoryDataHandlersV1 AddHistoryDataHandlersV1
func AddHistoryDataHandlersV1(group *gin.RouterGroup) {
	group.GET("/history/day_kbar/:stock/:start_date/:interval", GetKbarData)
}

// GetKbarData GetKbarData
// @Summary GetKbarData V1
// @tags HistoryData
// @accept json
// @produce json
// @param stock path string true "stock"
// @param start_date path string true "start_date"
// @param interval path string true "interval"
// @success 200 {object} []dbagent.HistoryKbar
// @failure 500 {object} ErrorResponse
// @Router /v1/history/day_kbar/{stock}/{start_date}/{interval} [get]
func GetKbarData(c *gin.Context) {
	stockNum := c.Param("stock")
	if target := cache.GetCache().GetTargetByStockNum(stockNum); target == nil {
		tmp := &dbagent.Target{
			Stock:       cache.GetCache().GetStock(stockNum),
			TradeDay:    cache.GetCache().GetTradeDay(),
			Rank:        0,
			Volume:      0,
			Subscribe:   false,
			RealTimeAdd: true,
		}
		// send to bus
		eventbus.Get().PublishTargets([]*dbagent.Target{tmp})
		// append to cache
		cache.GetCache().AppendTargets([]*dbagent.Target{tmp})
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	interval, err := strconv.Atoi(c.Param("interval"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	startDate := c.Param("start_date")
	startDateTime, err := time.Parse(global.ShortTimeLayout, startDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	var retryTimes int
	var result []dbagent.HistoryKbar
	for i := 0; i < interval; i++ {
		if retryTimes >= 300 {
			c.JSON(http.StatusInternalServerError, nil)
			return
		}
		tmp := cache.GetCache().GetStockHistoryDayKbar(startDateTime.AddDate(0, 0, -i).Format(global.ShortTimeLayout), stockNum)
		if i == 0 && tmp == nil {
			retryTimes++
			time.Sleep(100 * time.Millisecond)
			continue
		}
		if tmp != nil {
			result = append(result, *tmp)
		}
	}
	c.JSON(http.StatusOK, result)
}
