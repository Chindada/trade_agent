// Package routers package routers
package routers

import (
	"net/http"

	"trade_agent/pkg/config"

	"github.com/gin-gonic/gin"
)

// AddConfigHandlersV1 AddConfigHandlersV1
func AddConfigHandlersV1(group *gin.RouterGroup) {
	group.GET("/config", GetAllConfig)
}

// GetAllConfig GetAllConfig
// @Summary GetAllConfig V1
// @tags Config
// @accept json
// @produce json
// @success 200 {object} config.Config
// @Router /v1/config [get]
func GetAllConfig(c *gin.Context) {
	c.JSON(http.StatusOK, config.GetAllConfig())
}
