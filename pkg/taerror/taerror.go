// Package taerror package taerror
package taerror

import "fmt"

// TAError TAError
type TAError struct {
	ErrorCode int64  `json:"error_code,omitempty" yaml:"error_code"`
	Message   string `json:"message,omitempty" yaml:"message"`
	Err       error  `json:"err,omitempty" yaml:"err"`
}

// Error Error
func (c *TAError) Error() string {
	return c.Message
}

// ErrProtoFormatWrong ErrProtoFormatWrong
func ErrProtoFormatWrong(message []byte, err error) *TAError {
	return &TAError{
		ErrorCode: -101,
		Message:   fmt.Sprintf("Format Wrong: %s", string(message)),
		Err:       err,
	}
}
