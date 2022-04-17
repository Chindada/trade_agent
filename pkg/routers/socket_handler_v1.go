package routers

import (
	"trade_agent/pkg/websocket"

	"github.com/gin-gonic/gin"
)

// AddSocketHandlersV1 AddSocketHandlersV1
func AddSocketHandlersV1(group *gin.RouterGroup) {
	group.GET("/ws", SocketHandler)
}

// SocketHandler SocketHandler
func SocketHandler(c *gin.Context) {
	websocket.Run(c)
}
