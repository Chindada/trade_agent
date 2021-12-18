// Package mqhandler package mqhandler
package mqhandler

import (
	"sync"
	"trade_agent/pkg/config"
	"trade_agent/pkg/log"

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
	once.Do(initMQHandler)
	log.Get().Info("Initial MQHandler")
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
	var err error
	var conf config.Config
	conf, err = config.Get()
	if err != nil {
		log.Get().Panic(err)
	}
	var newClient mqtt.Client
	newClient, err = getMQClient(conf.GetMQConfig())
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
	c.lock.Lock()
	c.callbackMap[string(body.MQTopic)] = body.Callback
	if body.Once {
		c.onceMap[string(body.MQTopic)] = true
	}
	c.lock.Unlock()
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
	log.Get().Warnf("UnSubscribe %s", topic)
	return nil
}
