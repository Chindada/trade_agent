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

// addSwagger addSwagger
// @title Trade Agent
// @version 2.0.0
// @description API docs for Trade Agent
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath /trade-agent
func addSwagger(router *gin.Engine) {
	serverConf := config.GetServerConfig()
	var hostWithPort string
	if global.Get().GetIsDevelopment() {
		hostWithPort = fmt.Sprintf("%s:%s", utils.GetHostIP(), serverConf.HTTPPort)
	} else {
		hostWithPort = fmt.Sprintf("trade-agent.tocraw.com:%s", serverConf.HTTPPort)
	}
	docs.SwaggerInfo.Host = hostWithPort
	url := ginSwagger.URL("/swagger/doc.json")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
}
