// Package routers package routers
package routers

import (
	"net/http"

	"trade_agent/pkg/cache"

	"github.com/gin-gonic/gin"
)

// AddCacheHandlersV1 AddCacheHandlersV1
func AddCacheHandlersV1(group *gin.RouterGroup) {
	group.GET("/cache", GetAllCacheType)
	group.GET("/cache/:key", GetAllCacheDataByType)
}

// GetAllCacheType GetAllCacheType
// @Summary GetTradeDayTargets V1
// @tags Cache
// @accept json
// @produce json
// @success 200 {object} []string
// @Router /v1/cache [get]
func GetAllCacheType(c *gin.Context) {
	typeArr := cache.GetCache().GetAllCacheType()
	c.JSON(http.StatusOK, typeArr)
}

// GetAllCacheDataByType GetAllCacheDataByType
// @Summary GetTradeDayTargets V1
// @tags Cache
// @accept json
// @produce json
// @success 200 {object} interface{}
// @param key path string true "key"
// @Router /v1/cache/{key} [get]
func GetAllCacheDataByType(c *gin.Context) {
	key := c.Param("key")
	data := cache.GetCache().GetAllCacheByType(key)
	c.JSON(http.StatusOK, data)
}
