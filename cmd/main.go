package main

import (
	"trade_agent/pkg/log"
	"trade_agent/pkg/mqhandler"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	// dao.InitDatabase()
	// tasks.InitTasks()
	done := make(chan struct{})
	mqHandler, err := mqhandler.NewMQHandler()
	if err != nil {
		log.Get().Panic(err)
	}
	mqHandler.AddCallback(mqhandler.TradeRecordResponse, tradeRecordCallback(), false)
	mqHandler.Sub(mqhandler.TradeRecordResponse)
	<-done
}

func tradeRecordCallback() func(mqtt.Client, mqtt.Message) {
	return func(c mqtt.Client, m mqtt.Message) {
		log.Get().Warnf("topic: %s, msg: %s", m.Topic(), m.Payload())
	}
}
