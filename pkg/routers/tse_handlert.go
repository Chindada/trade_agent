// Package routers package routers
package routers

import (
	"net/http"

	"trade_agent/pkg/cache"

	"github.com/gin-gonic/gin"
)

// AddTSEHandlersV1 AddTSEHandlersV1
func AddTSEHandlersV1(group *gin.RouterGroup) {
	group.GET("/tse/real-time", GetRealTimeTSE)
}

// GetRealTimeTSE GetRealTimeTSE
// @Summary GetRealTimeTSE V1
// @tags TSE
// @accept json
// @produce json
// @success 200 {object} dbagent.TSESnapShot
// @Router /v1/tse/real-time [get]
func GetRealTimeTSE(c *gin.Context) {
	c.JSON(http.StatusOK, cache.GetCache().GetTSESnapshot())
}
