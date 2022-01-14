// Package taerror package taerror
package taerror

import (
	"fmt"
	"trade_agent/pkg/log"
)

// TAError TAError
type TAError struct {
	ErrorCode int   `json:"error_code,omitempty" yaml:"error_code"`
	Err       error `json:"err,omitempty" yaml:"err"`
}

// Error Error
func (c *TAError) Error() string {
	return c.Err.Error()
}

// ErrProtoFormatWrong ErrProtoFormatWrong
func ErrProtoFormatWrong(message []byte, err error) *TAError {
	code := -101
	protoErr := fmt.Errorf("format wrong: %s", string(message))
	log.Get().Errorf("ErrorCode: %d, Error: %s, %s", code, protoErr.Error(), err.Error())
	return &TAError{
		ErrorCode: code,
		Err:       protoErr,
	}
}
