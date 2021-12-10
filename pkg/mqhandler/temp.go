// Package mqhandler package mqhandler
package mqhandler

// func tradeRecordCallback() MQCallback {
// 	return func(c mqtt.Client, m mqtt.Message) {
// 		log.Get().Warnf("topic: %s, msg: %s", m.Topic(), m.Payload())
// 	}
// }

// func pubTopic(handler *MQHandler) {
// 	for {
// 		_ = handler.Pub(TradeRecordResponse, "123")
// 		time.Sleep(500 * time.Millisecond)
// 	}
// }

// // SubAndPub SubAndPub
// func SubAndPub() {
// 	handler := Get()
// 	handler.AddCallbackByTopic(TradeRecordResponse, tradeRecordCallback())
// 	_ = handler.Sub(TradeRecordResponse)
// 	pubTopic(handler)
// }
