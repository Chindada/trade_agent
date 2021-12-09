// Package mqhandler package mqhandler
package mqhandler

import (
	"sync"
	"trade_agent/pkg/config"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// MQHandler MQHandler
type MQHandler struct {
	lock        sync.Mutex
	client      mqtt.Client
	callbackMap map[string]topicCallback
}

type topicCallback struct {
	callBack func(mqtt.Client, mqtt.Message)
	once     bool
}

// NewMQHandler NewMQHandler
func NewMQHandler() (handler *MQHandler, err error) {
	var conf config.Config
	conf, err = config.Get()
	if err != nil {
		return nil, err
	}
	var newClient mqtt.Client
	newClient, err = getMQClient(conf.GetMQConfig())
	if err != nil {
		return nil, err
	}
	callbackMap := make(map[string]topicCallback)
	handler = &MQHandler{
		lock:        sync.Mutex{},
		client:      newClient,
		callbackMap: callbackMap,
	}
	return handler, err
}

// AddCallback AddCallback
func (c *MQHandler) AddCallback(topic Topic, callback func(mqtt.Client, mqtt.Message), once bool) {
	defer c.lock.Unlock()
	c.lock.Lock()
	tmp := topicCallback{
		callBack: callback,
		once:     once,
	}
	c.callbackMap[string(topic)] = tmp
}

// Pub Pub
func (c *MQHandler) Pub(topic Topic, msg interface{}) {
	c.client.Publish(string(topic), 0, false, msg)
}

// Sub Sub
func (c *MQHandler) Sub(topic Topic) {
	c.client.Subscribe(string(topic), 0, c.onMessage(topic))
}

// UnSub UnSub
func (c *MQHandler) UnSub(topics []Topic) {
	var tmp []string
	for _, v := range topics {
		tmp = append(tmp, string(v))
	}
	c.client.Unsubscribe(tmp...)
}

func (c *MQHandler) onMessage(topic Topic) func(mqtt.Client, mqtt.Message) {
	defer c.lock.Unlock()
	c.lock.Lock()
	callback := c.callbackMap[string(topic)]
	if callback.once {
		tmp := []Topic{topic}
		c.UnSub(tmp)
	}
	return callback.callBack
}
