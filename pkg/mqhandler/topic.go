// Package mqhandler package mqhandler
package mqhandler

// MQSubBody MQSubBody
type MQSubBody struct {
	Topic
	Once     bool
	Callback MQCallback
}

// Topic Topic
type Topic string

// TopicTradeRecord TopicTradeRecord
func TopicTradeRecord() Topic {
	return "internal/trade_record"
}
