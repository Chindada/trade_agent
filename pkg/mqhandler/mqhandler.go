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

// MQHandler MQHandler
type MQHandler struct {
	lock        sync.RWMutex
	client      mqtt.Client
	callbackMap map[string]MQCallback
}

// type topicCallback struct {
// 	callBack MQCallback
// 	once     bool
// }

// MQCallback MQCallback
type MQCallback func(mqtt.Client, mqtt.Message)

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
		log.Get().Panic("MQTT connect fail")
	}
	callbackMap := make(map[string]MQCallback)
	globalHandler = &MQHandler{
		lock:        sync.RWMutex{},
		client:      newClient,
		callbackMap: callbackMap,
	}
}

// AddCallbackByTopic AddCallbackByTopic
func (c *MQHandler) AddCallbackByTopic(topic Topic, cb MQCallback) {
	defer c.lock.Unlock()
	c.lock.Lock()
	c.callbackMap[string(topic)] = cb
}

// Pub Pub
func (c *MQHandler) Pub(topic Topic, msg interface{}) error {
	token := c.client.Publish(string(topic), 0, false, msg)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

// Sub Sub
func (c *MQHandler) Sub(topic Topic) error {
	token := c.client.Subscribe(string(topic), 0, c.onMessage(topic))
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

// UnSub UnSub
func (c *MQHandler) UnSub(topic Topic) error {
	var tmp []string
	tmp = append(tmp, string(topic))

	token := c.client.Unsubscribe(tmp...)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (c *MQHandler) onMessage(topic Topic) func(mqtt.Client, mqtt.Message) {
	defer c.lock.RUnlock()
	c.lock.RLock()
	callback := c.callbackMap[string(topic)]
	return callback
}
