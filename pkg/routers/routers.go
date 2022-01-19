// Package routers package routers
package routers

import (
	"net/http"
	"trade_agent/docs"
	"trade_agent/pkg/config"
	"trade_agent/pkg/log"
	"trade_agent/pkg/routers/handlers/balance"
	"trade_agent/pkg/routers/handlers/targets"

	"github.com/gin-gonic/gin"
)

// ServeHTTP ServeHTTP
func ServeHTTP() {
	go func() {
		serverConf := config.GetServerConfig()
		gin.SetMode(serverConf.RunMode)
		g := gin.New()
		g.Use(corsMiddleware())
		g.Use(gin.Recovery())
		err := g.SetTrustedProxies(nil)
		if err != nil {
			log.Get().Panic(err)
		}
		addSwagger(g)
		initRouters(g)
		log.Get().Infof("HTTP Server On %s", docs.SwaggerInfo.Host)
		if err := g.Run(":" + serverConf.HTTPPort); err != nil {
			log.Get().Panic(err)
		}
	}()
}

func initRouters(router *gin.Engine) {
	mainRoute := router.Group("trade-agent")
	balance.AddHandlers(mainRoute)
	targets.AddHandlers(mainRoute)
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "*")
		c.Header("Access-Control-Allow-Methods", "GET, OPTIONS, POST, PUT, DELETE")
		c.Set("content-type", "application/json")
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, nil)
		}
		c.Next()
	}
}
