// Package sinopacerr package sinopacerr
package sinopacerr

import (
	"errors"
	"trade_agent/pkg/log"
)

// SinopacError SinopacError
type SinopacError struct {
	ErrorCode int   `json:"error_code,omitempty" yaml:"error_code"`
	Err       error `json:"err,omitempty" yaml:"err"`
}

// Error Error
func (c *SinopacError) Error() string {
	return c.Err.Error()
}

// ErrQuotaIsNotEnough ErrQuotaIsNotEnough
func ErrQuotaIsNotEnough() *SinopacError {
	code := -501
	err := errors.New(QuotaIsNotEnough)
	log.Get().Errorf("ErrorCode: %d, Error: %s", code, err.Error())
	return &SinopacError{
		ErrorCode: code,
		Err:       err,
	}
}
