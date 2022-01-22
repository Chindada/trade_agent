// Package routers package routers
package routers

import (
	"net/http"
	"trade_agent/pkg/config"

	"github.com/gin-gonic/gin"
)

// AddConfigHandlers AddConfigHandlers
func AddConfigHandlers(group *gin.RouterGroup) {
	group.GET("/config", GetAllConfig)
}

// GetAllConfig GetAllConfig
// @Summary GetAllConfig
// @tags Config
// @accept json
// @produce json
// @success 200 {object} config.Config
// @Router /config [get]
func GetAllConfig(c *gin.Context) {
	c.JSON(http.StatusOK, config.GetAllConfig())
}
