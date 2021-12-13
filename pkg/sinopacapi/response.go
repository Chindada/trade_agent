// Package sinopacapi package sinopacapi
package sinopacapi

// ResponseHealthStatus ResponseHealthStatus
type ResponseHealthStatus struct {
	Result      string `json:"result"`
	UpTimeMin   int64  `json:"up_time_min"`
	ServerToken string `json:"server_token"`
}

// OrderResponse OrderResponse
type OrderResponse struct {
	Status  string `json:"status"`
	OrderID string `json:"order_id"`
}

// ResponseCommon ResponseCommon
type ResponseCommon struct {
	Result string `json:"result"`
}
