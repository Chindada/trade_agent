// Package routers package routers
package routers

import (
	"net/http"
	"sort"
	"time"

	"trade_agent/global"
	"trade_agent/pkg/cache"
	"trade_agent/pkg/config"
	"trade_agent/pkg/dbagent"
	"trade_agent/pkg/eventbus"
	"trade_agent/pkg/log"
	"trade_agent/pkg/modules/tradeday"

	"github.com/gin-gonic/gin"
)

// AddTargetsHandlers AddTargetsHandlers
func AddTargetsHandlers(group *gin.RouterGroup) {
	group.GET("/targets", GetTradeDayTargets)
	group.POST("/targets", AddTargets)

	group.GET("/targets/quater", GetQuaterTargets)
}

// GetTradeDayTargets GetTradeDayTargets
// @Summary GetTradeDayTargets
// @tags Targets
// @accept json
// @produce json
// @success 200 {object} []dbagent.Target
// @failure 500 {object} ErrorResponse
// @Router /targets [get]
func GetTradeDayTargets(c *gin.Context) {
	targets := cache.GetCache().GetTargets()
	c.JSON(http.StatusOK, targets)
}

// AddTargets AddTargets
// @Summary GetTradeDayTargets
// @tags Targets
// @accept json
// @produce json
// @param price_range header string true "price_range"
// @success 200
// @failure 500 {object} ErrorResponse
// @Router /targets [post]
func AddTargets(c *gin.Context) {
	originalTargets := cache.GetCache().GetTargets()
	originalMap := make(map[string]bool)
	for _, v := range originalTargets {
		originalMap[v.Stock.Number] = true
	}
	var targets []*dbagent.Target
	priceRange := c.Request.Header.Get("price_range")

	var err error
	switch priceRange {
	case "1":
		targets, err = queryAllStockByMinMax(0, 10, originalMap)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
	case "2":
		targets, err = queryAllStockByMinMax(10, 50, originalMap)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
	case "3":
		targets, err = queryAllStockByMinMax(50, 100, originalMap)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
	case "4":
		targets, err = queryAllStockByMinMax(100, 500, originalMap)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
	}

	if len(targets) != 0 {
		// send to bus
		eventbus.Get().PublishTargets(targets)
		// append to cache
		cache.GetCache().AppendTargets(targets)
	}

	c.JSON(http.StatusOK, gin.H{
		"total_add": len(targets),
	})
}

func queryAllStockByMinMax(min, max float64, originalMap map[string]bool) ([]*dbagent.Target, error) {
	var tmp []*dbagent.Target
	allStock, err := dbagent.Get().GetAllStockMap()
	if err != nil {
		return []*dbagent.Target{}, err
	}
	conf := config.GetTargetCondConfig()
	for stock := range allStock {
		if cache.GetCache().GetStockVolume(stock) < conf.LimitVolume {
			continue
		}
		if allStock[stock].LastClose >= min && allStock[stock].LastClose < max && !originalMap[stock] {
			tmp = append(tmp, &dbagent.Target{
				Stock:       allStock[stock],
				TradeDay:    cache.GetCache().GetTradeDay(),
				Rank:        -1,
				Volume:      -1,
				Subscribe:   false,
				RealTimeAdd: true,
			})
		}
	}
	return tmp, nil
}

// GetQuaterTargets GetQuaterTargets
// @Summary GetTradeDayTargets
// @tags Targets
// @accept json
// @produce json
// @success 200 {object} []QuaterMAResponse{}
// @failure 500 {object} ErrorResponse
// @Router /targets/quater [get]
func GetQuaterTargets(c *gin.Context) {
	targets := cache.GetCache().GetTargets()
	belowQuater := make(map[time.Time][]dbagent.Stock)
	result := []QuaterMAResponse{}
	for _, t := range targets {
		tmp := *t
		maArr, err := dbagent.Get().GetAllQuaterMAByStockID(int64(tmp.Stock.ID))
		if err != nil {
			log.Get().Error(err)
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		for _, ma := range maArr {
			if close := cache.GetCache().GetHistoryClose(ma.Stock.Number, ma.CalendarDate.Date); close != 0 && close-ma.QuaterMA < 0 {
				nextTradeDay, err := tradeday.GetAbsNextTradeDayTime(ma.CalendarDate.Date)
				if err != nil {
					log.Get().Error(err)
					c.JSON(http.StatusInternalServerError, err)
					return
				}
				if nextOpen := cache.GetCache().GetHistoryOpen(ma.Stock.Number, nextTradeDay); nextOpen != 0 && nextOpen-ma.QuaterMA > 0 {
					belowQuater[ma.CalendarDate.Date] = append(belowQuater[ma.CalendarDate.Date], *tmp.Stock)
				}
			}
		}
	}

	dateArr := []time.Time{}
	for date := range belowQuater {
		dateArr = append(dateArr, date)
	}

	sort.Slice(dateArr, func(i, j int) bool {
		return dateArr[i].After(dateArr[j])
	})

	for _, date := range dateArr {
		result = append(result, QuaterMAResponse{
			Date:   date.Format(global.ShortTimeLayout),
			Stocks: belowQuater[date],
		})
	}
	c.JSON(http.StatusOK, result)
}
