// Package dbagent package dbagent
package dbagent

import "time"

// TSESnapShot TSESnapShot
type TSESnapShot struct {
	Stock    *Stock    `json:"stock,omitempty" yaml:"stock"`
	TickTime time.Time `json:"tick_time,omitempty" yaml:"tick_time"`

	Open            float64 `json:"open,omitempty" yaml:"open"`
	High            float64 `json:"high,omitempty" yaml:"high"`
	Low             float64 `json:"low,omitempty" yaml:"low"`
	Close           float64 `json:"close,omitempty" yaml:"close"`
	TickType        string  `json:"tick_type,omitempty" yaml:"tick_type"`
	PriceChg        float64 `json:"price_chg,omitempty" yaml:"price_chg"`
	PctChg          float64 `json:"pct_chg,omitempty" yaml:"pct_chg"`
	ChgType         string  `json:"chg_type,omitempty" yaml:"chg_type"`
	Volume          int64   `json:"volume,omitempty" yaml:"volume"`
	VolumeSum       int64   `json:"volume_sum,omitempty" yaml:"volume_sum"`
	Amount          int64   `json:"amount,omitempty" yaml:"amount"`
	AmountSum       int64   `json:"amount_sum,omitempty" yaml:"amount_sum"`
	YesterdayVolume float64 `json:"yesterday_volume,omitempty" yaml:"yesterday_volume"`
	VolumeRatio     float64 `json:"volume_ratio,omitempty" yaml:"volume_ratio"`
}
