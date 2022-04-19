// Package routers package routers
package routers

// ErrorResponse ErrorResponse
type ErrorResponse struct {
	Response   string      `json:"response,omitempty" yaml:"response"`
	Attachment interface{} `json:"attachment,omitempty" yaml:"attachment"`
}
