// Package mqhandler package mqhandler
package mqhandler

import (
	"sync"
	"time"

	"trade_agent/pkg/config"
	"trade_agent/pkg/log"
	"trade_agent/pkg/utils"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var (
	globalHandler *MQHandler
	once          sync.Once
)

// MQCallback MQCallback
type MQCallback func(MQMessage)

// MQMessage MQMessage
type MQMessage mqtt.Message

// MQSubBody MQSubBody
type MQSubBody struct {
	MQTopic  MQTopic
	Once     bool
	Callback MQCallback
}

// MQHandler MQHandler
type MQHandler struct {
	lock        sync.RWMutex
	client      mqtt.Client
	callbackMap map[string]MQCallback
	onceMap     map[string]bool
}

// InitMQHandler InitMQHandler
func InitMQHandler() {
	log.Get().Info("Initial MQHandler")

	once.Do(initMQHandler)
}

// Get Get
func Get() *MQHandler {
	if globalHandler != nil {
		return globalHandler
	}
	log.Get().Panic("MQHandler was not initailized")
	return nil
}

// initMQHandler initMQHandler
func initMQHandler() {
	if globalHandler != nil {
		return
	}

	mqConf := config.GetMQTTConfig()
	for {
		if utils.CheckPortIsOpen(mqConf.Host, mqConf.Port) {
			break
		}
		time.Sleep(time.Second)
	}

	newClient, err := getMQClient(mqConf)
	if err != nil {
		log.Get().Panic(err)
	}

	if newClient == nil {
		log.Get().Panic("MQTT Connect Fail")
	}

	callbackMap := make(map[string]MQCallback)
	onceMap := make(map[string]bool)
	globalHandler = &MQHandler{
		lock:        sync.RWMutex{},
		client:      newClient,
		callbackMap: callbackMap,
		onceMap:     onceMap,
	}
}

// Sub Sub
func (c *MQHandler) Sub(body MQSubBody) error {
	var already bool
	c.lock.Lock()
	// check if already in callback map
	if c.callbackMap[string(body.MQTopic)] != nil {
		already = true
	} else {
		c.callbackMap[string(body.MQTopic)] = body.Callback
	}
	// overwrite once map
	if body.Once {
		c.onceMap[string(body.MQTopic)] = true
	} else {
		delete(c.onceMap, string(body.MQTopic))
	}
	c.lock.Unlock()
	if already {
		return nil
	}
	go func() {
		token := c.client.Subscribe(string(body.MQTopic), 2, c.onMessage)
		if token.Wait() && token.Error() != nil {
			log.Get().Error(token.Error())
		}
	}()
	return nil
}

// Pub Pub
func (c *MQHandler) Pub(topic MQTopic, msg interface{}) error {
	token := c.client.Publish(string(topic), 2, false, msg)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (c *MQHandler) onMessage(_ mqtt.Client, m mqtt.Message) {
	c.lock.RLock()
	if f := c.callbackMap[m.Topic()]; f != nil {
		go f(m)
	}
	once := c.onceMap[m.Topic()]
	c.lock.RUnlock()
	if once {
		c.lock.Lock()
		delete(c.callbackMap, m.Topic())
		delete(c.onceMap, m.Topic())
		c.lock.Unlock()
		err := c.UnSub(m.Topic())
		if err != nil {
			log.Get().Error(err)
		}
	}
}

// UnSub UnSub
func (c *MQHandler) UnSub(topic string) error {
	token := c.client.Unsubscribe(topic)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	log.Get().WithFields(map[string]interface{}{
		"Topic": topic,
	}).Warn("UnSubscribe")
	return nil
}
