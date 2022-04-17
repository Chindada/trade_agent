// Package routers package routers
package routers

import (
	"fmt"
	"net/http"

	"trade_agent/global"
	"trade_agent/pkg/config"
	"trade_agent/pkg/log"

	"github.com/gin-gonic/gin"
)

// APIVersion APIVersion
const APIVersion string = "v1"

var loginFunc gin.HandlerFunc

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

		authMiddleware := AuthMiddleware(g)
		addSwagger(g)

		public := g.Group(fmt.Sprintf("%s/%s", basePath, APIVersion))
		loginFunc = authMiddleware.LoginHandler
		public.POST("/login", loginHandler())

		private := g.Group(fmt.Sprintf("%s/%s", basePath, APIVersion))
		private.Use(authMiddleware.MiddlewareFunc())

		AddBalanceHandlersV1(public)
		AddTargetsHandlersV1(public)
		AddOrderHandlersV1(public)
		AddConfigHandlersV1(public)
		AddTSEHandlersV1(public)
		AddSocketHandlersV1(public)

		// log.Get().Infof("HTTP Server On %s", docs.SwaggerInfo.Host)
		listenPath := fmt.Sprintf(":%s", serverConf.HTTPPort)
		if global.Get().GetIsDevelopment() {
			AddCacheHandlersV1(public)
			if err := g.Run(listenPath); err != nil {
				log.Get().Panic(err)
			}
		} else {
			if err := g.RunTLS(listenPath, serverConf.CertPath, serverConf.KeyPath); err != nil {
				log.Get().Panic(err)
			}
		}
	}()
}

// loginHandler loginHandler
// @tags Login V1
// @accept json
// @produce json
// @param body body dbagent.Login{} true "Body"
// @success 200
// @Router /v1/login [post]
func loginHandler() gin.HandlerFunc {
	return loginFunc
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
