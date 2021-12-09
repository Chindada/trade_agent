// Package routers package routers
package routers

import (
	"trade_agent/docs"
	"trade_agent/pkg/config"
	"trade_agent/pkg/log"

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
	conf, err := config.Get()
	if err != nil {
		log.Get().Panic(err)
	}
	serverConf := conf.GetServerConfig()
	docs.SwaggerInfo.Host = "trade-agent.tocraw.com:" + serverConf.HTTPPort
	url := ginSwagger.URL("/swagger/doc.json")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
}
