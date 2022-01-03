// Package handlers package handlers
package handlers

// ErrorResponse ErrorResponse
type ErrorResponse struct {
	Response   string      `json:"response,omitempty" yaml:"response"`
	Attachment interface{} `json:"attachment,omitempty" yaml:"attachment"`
}
