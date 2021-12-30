// Package history package history
package history

import (
	"errors"
	"trade_agent/pkg/utils"

	"github.com/markcheno/go-talib"
)

func generareMAByCloseArr(closeArr []float64) (lastMa float64, err error) {
	maArr := talib.Ma(closeArr, len(closeArr), talib.SMA)
	if len(maArr) == 0 {
		return 0, errors.New("no ma")
	}
	return maArr[len(maArr)-1], err
}

func getBiasRateByCloseArr(closeArr []float64) (biasRate float64, err error) {
	var ma float64
	ma, err = generareMAByCloseArr(closeArr)
	if err != nil {
		return biasRate, err
	}
	biasRate = utils.Round(100*(closeArr[len(closeArr)-1]-ma)/ma, 2)
	return biasRate, err
}
