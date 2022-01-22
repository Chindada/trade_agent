// Package routers package routers
package routers

import (
	"net/http"
	"trade_agent/pkg/cache"

	"github.com/gin-gonic/gin"
)

// AddTargetsHandlers AddTargetsHandlers
func AddTargetsHandlers(group *gin.RouterGroup) {
	group.GET("/targets", GetTradeDayTargets)
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
