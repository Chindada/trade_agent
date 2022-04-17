// Package websocket package websocket
package websocket

import (
	"encoding/json"
	"net/http"
	"os"
	"sync"

	"trade_agent/pkg/log"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type wsClient struct {
	PickStock []string    `json:"pick_stock" yaml:"pick_stock"`
	Msg       interface{} `json:"msg" yaml:"msg"`

	lock sync.Mutex      `json:"-" yaml:"lock"`
	Conn *websocket.Conn `json:"conn" yaml:"conn"`
}

type msg struct {
	Data interface{} `json:"data" yaml:"data"`
}

var (
	upGrader    = websocket.Upgrader{}
	offline     = make(chan *websocket.Conn)
	msgChan     = make(chan *wsClient)
	activeUser  = make(map[string]*wsClient) // key: remoteAddr
	userMapLock sync.Mutex
)

// Run Run
func Run(gin *gin.Context) {
	upGrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	c, err := upGrader.Upgrade(gin.Writer, gin.Request, nil)
	if err != nil {
		log.Get().Errorf("WS Upgrade: %s", err)
		return
	}

	defer func() {
		if err := c.Close(); err != nil {
			log.Get().Errorf("Close error: %s", err)
		}
	}()

	stuck := make(chan os.Signal)
	go read(c)
	go sendMsg()
	go kickOut()

	go processPickStock()
	<-stuck
}

func read(c *websocket.Conn) {
	defer func() {
		if err := recover(); err != nil {
			log.Get().Errorf("Read error: %s", err)
		}
	}()

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			offline <- c
			return
		}

		serveMsgStr := message
		if string(serveMsgStr) == "ping" {
			_ = c.WriteMessage(websocket.TextMessage, []byte("pong"))
			continue
		}

		clientMsg := msg{}
		_ = json.Unmarshal(message, &clientMsg)
		if clientMsg.Data != nil {
			var pickStock []string
			if _, ok := clientMsg.Data.(map[string]interface{})["pick_stock_list"].([]interface{}); ok {
				for _, v := range clientMsg.Data.(map[string]interface{})["pick_stock_list"].([]interface{}) {
					pickStock = append(pickStock, v.(string))
				}
			}
			userMapLock.Lock()
			activeUser[c.RemoteAddr().String()] = &wsClient{
				PickStock: pickStock,
				Msg:       nil,
				lock:      sync.Mutex{},
				Conn:      c,
			}
			userMapLock.Unlock()
			// _ = c.WriteMessage(websocket.TextMessage, "success")
		}
	}
}

func sendMsg() {
	var user string
	defer func() {
		if err := recover(); err != nil {
			log.Get().Errorf("Write error: %s:%s", user, err)
		}
	}()
	for {
		msg := <-msgChan
		msg.lock.Lock()
		if msg.Conn != nil && msg.Msg != nil {
			user = msg.Conn.RemoteAddr().String()
			serveMsgStr, _ := json.Marshal(msg.Msg)
			// fmt.Println(msg.Msg, user)
			go func() {
				_ = msg.Conn.WriteMessage(websocket.TextMessage, serveMsgStr)
			}()
		}
		msg.lock.Unlock()
	}
}

func kickOut() {
	for {
		out := <-offline
		if _, ok := activeUser[out.RemoteAddr().String()]; ok {
			userMapLock.Lock()
			delete(activeUser, out.RemoteAddr().String())
			userMapLock.Unlock()
		}
	}
}
