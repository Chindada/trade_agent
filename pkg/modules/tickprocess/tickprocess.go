// Package tickprocess package tickprocess
package tickprocess

import (
	"trade_agent/pkg/log"
	"trade_agent/pkg/mqhandler"
)

// InitTickProcess InitTickProcess
func InitTickProcess() {
	handler := mqhandler.Get()
	// real time
	realTimeSubBody := mqhandler.MQSubBody{
		MQTopic:  mqhandler.TopicRealTimeTick(),
		Once:     false,
		Callback: realTimeTickCallback,
	}
	err := handler.Sub(realTimeSubBody)
	if err != nil {
		log.Get().Panic(err)
	}

	// history tick
	historyTickSubBody := mqhandler.MQSubBody{
		MQTopic:  mqhandler.TopicHistoryTick(),
		Once:     false,
		Callback: historyTickCallback,
	}
	err = handler.Sub(historyTickSubBody)
	if err != nil {
		log.Get().Panic(err)
	}
}
