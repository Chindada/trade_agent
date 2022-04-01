// Package routers package routers
package routers

import "trade_agent/pkg/dbagent"

// ErrorResponse ErrorResponse
type ErrorResponse struct {
	Response   string      `json:"response,omitempty" yaml:"response"`
	Attachment interface{} `json:"attachment,omitempty" yaml:"attachment"`
}

// QuaterMAResponse QuaterMAResponse
type QuaterMAResponse struct {
	Date   string          `json:"date"`
	Stocks []dbagent.Stock `json:"stocks"`
}
