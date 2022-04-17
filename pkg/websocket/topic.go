package websocket

// Topic Topic
type Topic string

func (c Topic) String() string {
	return string(c)
}

// WSTopicPickStock WSTopicPickStock
const WSTopicPickStock Topic = "pick_stock"
