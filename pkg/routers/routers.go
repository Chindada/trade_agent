// Package routers package routers
package routers

import (
	"trade_agent/pkg/config"
	"trade_agent/pkg/log"
	"trade_agent/pkg/utils"

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
		log.Get().Infof("HTTP Server serve on http://%s:%s/", utils.GetHostIP(), serverConf.HTTPPort)
		if err := g.Run(":" + serverConf.HTTPPort); err != nil {
			log.Get().Panic(err)
		}
	}()
}

func initRouters(router *gin.Engine) {
	// mainRoute := router.Group("trade-agent")
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "*")
		c.Header("Access-Control-Allow-Methods", "GET, OPTIONS, POST, PUT, DELETE")
		c.Set("content-type", "application/json")
		c.Next()
	}
}
