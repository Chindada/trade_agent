// Package targets package targets
package targets

import (
	"net/http"
	"trade_agent/pkg/cache"

	"github.com/gin-gonic/gin"
)

// AddHandlers AddHandlers
func AddHandlers(group *gin.RouterGroup) {
	group.GET("/targets", GetTradeDayTargets)
}

// GetTradeDayTargets GetTradeDayTargets
// @Summary GetTradeDayTargets
// @tags Targets
// @accept json
// @produce json
// @success 200 {object} []dbagent.Target
// @failure 500 {object} handlers.ErrorResponse
// @Router /targets [get]
func GetTradeDayTargets(c *gin.Context) {
	targets := cache.GetCache().GetTargets()
	c.JSON(http.StatusOK, targets)
}
