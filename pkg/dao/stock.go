// Package dao package dao
package dao

import "gorm.io/gorm"

// Stock Stock
type Stock struct {
	gorm.Model `json:"-" swaggerignore:"true"`

	Number             string  `json:"number" gorm:"column:number;uniqueIndex;index:idx_stock"`
	Name               string  `json:"name" gorm:"column:name"`
	Type               string  `json:"type" gorm:"column:type"`
	Category           string  `json:"category" gorm:"column:category"`
	DayTrade           bool    `json:"day_trade" gorm:"column:day_trade;index:idx_stock"`
	LastClose          float64 `json:"last_close" gorm:"column:last_close"`
	LastVolume         int64   `json:"last_volume" gorm:"column:last_volume"`
	LastCloseChangePct float64 `json:"last_close_change_pct" gorm:"column:last_close_change_pct"`
}
