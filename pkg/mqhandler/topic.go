// Package mqhandler package mqhandler
package mqhandler

// MQTopic MQTopic
type MQTopic string

// TopicStockDetail TopicStockDetail
func TopicStockDetail() MQTopic {
	return "internal/stock_detail"
}

// TopicSnapshotAll TopicSnapshotAll
func TopicSnapshotAll() MQTopic {
	return "internal/snapshot_all"
}

// TopicSnapshotTSE TopicSnapshotTSE
func TopicSnapshotTSE() MQTopic {
	return "internal/snapshot_tse"
}

// TopicHistoryTick TopicHistoryTick
func TopicHistoryTick() MQTopic {
	return "internal/history_tick"
}

// TopicHistoryTickTSE TopicHistoryTickTSE
func TopicHistoryTickTSE() MQTopic {
	return "internal/history_tick_tse"
}

// TopicHistoryKbar TopicHistoryKbar
func TopicHistoryKbar() MQTopic {
	return "internal/history_kbar"
}

// TopicHistoryKbarTSE TopicHistoryKbarTSE
func TopicHistoryKbarTSE() MQTopic {
	return "internal/history_kbar_tse"
}

// TopicVolumeRank TopicVolumeRank
func TopicVolumeRank() MQTopic {
	return "internal/volumerank"
}

// TopicOrderStatus TopicOrderStatus
func TopicOrderStatus() MQTopic {
	return "internal/order_status"
}

// TopicRealTimeTick TopicRealTimeTick
func TopicRealTimeTick() MQTopic {
	return "internal/realtime_tick"
}

// TopicRealTimeBidask TopicRealTimeBidask
func TopicRealTimeBidask() MQTopic {
	return "internal/realtime_bidask"
}

// TopicTradeEvent TopicTradeEvent
func TopicTradeEvent() MQTopic {
	return "internal/trade_event"
}

// TopicLastcount TopicLastcount
func TopicLastcount() MQTopic {
	return "internal/lastcount"
}

// TopicLastcountTSE TopicLastcountTSE
func TopicLastcountTSE() MQTopic {
	return "internal/lastcount_tse"
}

// TopicLastcountMultiDate TopicLastcountMultiDate
func TopicLastcountMultiDate() MQTopic {
	return "internal/lastcount_multi_date"
}
