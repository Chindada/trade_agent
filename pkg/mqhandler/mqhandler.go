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
type MQCallback func(mqtt.Client, mqtt.Message)

// MQHandler MQHandler
type MQHandler struct {
	lock        sync.Mutex
	client      mqtt.Client
	callbackMap map[string]MQCallback
}

// Get Get
func Get() *MQHandler {
	if globalHandler != nil {
		return globalHandler
	}
	once.Do(initMQHandler)
	return globalHandler
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
	globalHandler = &MQHandler{
		lock:        sync.Mutex{},
		client:      newClient,
		callbackMap: callbackMap,
	}
}

// Sub Sub
func (c *MQHandler) Sub(body MQSubBody) error {
	c.lock.Lock()
	c.callbackMap[string(body.Topic)] = body.Callback
	c.lock.Unlock()
	token := c.client.Subscribe(string(body.Topic), 2, c.onMessage(body))
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

// Pub Pub
func (c *MQHandler) Pub(topic Topic, msg interface{}) error {
	token := c.client.Publish(string(topic), 2, false, msg)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (c *MQHandler) onMessage(body MQSubBody) func(mqtt.Client, mqtt.Message) {
	defer c.lock.Unlock()
	c.lock.Lock()
	callback := c.callbackMap[string(body.Topic)]
	if callback == nil {
		err := c.UnSub(string(body.Topic))
		if err != nil {
			log.Get().Error(err)
		}
		log.Get().Infof("Finish %s Once Subscribe", body.Topic)
	}
	if body.Once {
		delete(c.callbackMap, string(body.Topic))
	}
	return callback
}

// UnSub UnSub
func (c *MQHandler) UnSub(topic string) error {
	token := c.client.Unsubscribe(topic)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}
