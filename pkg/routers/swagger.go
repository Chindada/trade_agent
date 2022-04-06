// Package routers package routers
package routers

import (
	"fmt"

	"trade_agent/docs"
	"trade_agent/global"
	"trade_agent/pkg/config"
	"trade_agent/pkg/utils"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	swaggerFiles "github.com/swaggo/gin-swagger/swaggerFiles"
)

var basePath string = "/trade-agent"

// @title Trade Agent
// @description API docs for Trade Agent

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
func addSwagger(router *gin.Engine) {
	serverConf := config.GetServerConfig()

	var hostWithPort string
	if global.Get().GetIsDevelopment() {
		hostWithPort = fmt.Sprintf("%s:%s", utils.GetHostIP(), serverConf.HTTPPort)
	} else {
		hostWithPort = fmt.Sprintf("trader.tocandraw.com:%s", serverConf.HTTPPort)
	}

	docs.SwaggerInfo.Host = hostWithPort
	docs.SwaggerInfo.BasePath = basePath
	docs.SwaggerInfo.Version = "v2.0.0"

	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
