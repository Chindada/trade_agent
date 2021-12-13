// Package mqhandler package mqhandler
package mqhandler

// MQSubBody MQSubBody
type MQSubBody struct {
	MQTopic
	Once     bool
	Callback MQCallback
}

// MQTopic MQTopic
type MQTopic string

// TopicTradeRecord TopicTradeRecord
func TopicTradeRecord() MQTopic {
	return "internal/trade_record"
}

// TopicStockDetail TopicStockDetail
func TopicStockDetail() MQTopic {
	return "internal/stock_detail"
}
