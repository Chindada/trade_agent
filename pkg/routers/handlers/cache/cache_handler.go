// Package cache package cache
package cache

import (
	"net/http"
	"trade_agent/pkg/cache"

	"github.com/gin-gonic/gin"
)

// AddHandlers AddHandlers
func AddHandlers(group *gin.RouterGroup) {
	group.GET("/cache", GetAllCacheType)
	group.GET("/cache/:key", GetAllCacheDataByType)
}

// GetAllCacheType GetAllCacheType
// @Summary GetTradeDayTargets
// @tags Cache
// @accept json
// @produce json
// @success 200 {object} []string
// @Router /cache [get]
func GetAllCacheType(c *gin.Context) {
	typeArr := cache.GetCache().GetAllCacheType()
	c.JSON(http.StatusOK, typeArr)
}

// GetAllCacheDataByType GetAllCacheDataByType
// @Summary GetTradeDayTargets
// @tags Cache
// @accept json
// @produce json
// @success 200 {object} interface{}
// @param key path string true "key"
// @Router /cache/{key} [get]
func GetAllCacheDataByType(c *gin.Context) {
	key := c.Param("key")
	data := cache.GetCache().GetAllCacheByType(key)
	c.JSON(http.StatusOK, data)
}
