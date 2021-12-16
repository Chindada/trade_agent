// Package utils package utils
package utils

import (
	"math"
	"strconv"
)

// StrToInt64 StrToInt64
func StrToInt64(input string) (ans int64, err error) {
	ans, err = strconv.ParseInt(input, 10, 64)
	if err != nil {
		return ans, err
	}
	return ans, err
}

// StrToFloat64 StrToFloat64
func StrToFloat64(input string) (ans float64, err error) {
	ans, err = strconv.ParseFloat(input, 64)
	if err != nil {
		return ans, err
	}
	return ans, err
}

// Round Round
func Round(val float64, precision int) float64 {
	p := math.Pow10(precision)
	return math.Floor(val*p+0.5) / p
}
