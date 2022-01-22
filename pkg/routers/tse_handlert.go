// Package routers package routers
package routers

import (
	"net/http"
	"trade_agent/pkg/cache"

	"github.com/gin-gonic/gin"
)

// AddTSEHandlers AddTSEHandlers
func AddTSEHandlers(group *gin.RouterGroup) {
	group.GET("/tse/real-time", GetRealTimeTSE)
}

// GetRealTimeTSE GetRealTimeTSE
// @Summary GetRealTimeTSE
// @tags TSE
// @accept json
// @produce json
// @success 200 {object} dbagent.TSESnapShot
// @Router /tse/real-time [get]
func GetRealTimeTSE(c *gin.Context) {
	c.JSON(http.StatusOK, cache.GetCache().GetTSESnapshot())
}
