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
	MQTopic  MQTopic    `json:"mq_topic,omitempty" yaml:"mq_topic"`
	Once     bool       `json:"once,omitempty" yaml:"once"`
	Callback MQCallback `json:"callback,omitempty" yaml:"callback"`
}

// MQHandler MQHandler
type MQHandler struct {
	lock        sync.Mutex
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
		lock:        sync.Mutex{},
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
	token := c.client.Subscribe(string(body.MQTopic), 2, c.onMessage)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
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

func (c *MQHandler) onMessage(mc mqtt.Client, m mqtt.Message) {
	defer c.lock.Unlock()
	c.lock.Lock()
	callback := c.callbackMap[m.Topic()]
	if c.onceMap[m.Topic()] {
		delete(c.callbackMap, m.Topic())
		delete(c.onceMap, m.Topic())
		err := c.UnSub(m.Topic())
		if err != nil {
			log.Get().Error(err)
		}
		log.Get().Warnf("UnSubscribe %s", m.Topic())
	}
	callback(m)
}

// UnSub UnSub
func (c *MQHandler) UnSub(topic string) error {
	token := c.client.Unsubscribe(topic)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}
