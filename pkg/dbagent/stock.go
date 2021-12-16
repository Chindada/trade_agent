// Package dbagent package dbagent
package dbagent

import "gorm.io/gorm"

// Stock Stock
type Stock struct {
	gorm.Model `json:"-" swaggerignore:"true"`

	Number             string  `json:"number,omitempty" yaml:"number" gorm:"column:number"`
	Name               string  `json:"name,omitempty" yaml:"name" gorm:"column:name"`
	Exchange           string  `json:"exchange,omitempty" yaml:"exchange" gorm:"column:exchange"`
	Category           string  `json:"category,omitempty" yaml:"category" gorm:"column:category"`
	DayTrade           bool    `json:"day_trade,omitempty" yaml:"day_trade" gorm:"column:day_trade"`
	LastClose          float64 `json:"last_close,omitempty" yaml:"last_close" gorm:"column:last_close"`
	LastVolume         int64   `json:"last_volume,omitempty" yaml:"last_volume" gorm:"column:last_volume"`
	LastCloseChangePct float64 `json:"last_close_change_pct,omitempty" yaml:"last_close_change_pct" gorm:"column:last_close_change_pct"`
}

// InsertStock InsertStock
func (c *DBAgent) InsertStock(stock *Stock) error {
	err := c.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&stock).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// InsertMultiStock InsertMultiStock
func (c *DBAgent) InsertMultiStock(stockArr []*Stock) (err error) {
	err = c.DB.Transaction(func(tx *gorm.DB) error {
		if err = tx.CreateInBatches(&stockArr, multiInsertBatchSize).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// GetAllStockMap GetAllStockMap
func (c *DBAgent) GetAllStockMap() (allStockMap map[string]*Stock, err error) {
	allStockMap = make(map[string]*Stock)
	var stockArr []*Stock
	err = c.DB.Model(&Stock{}).Not("name = ?", "").Order("number").Find(&stockArr).Error

	for _, v := range stockArr {
		allStockMap[v.Number] = v
	}
	return allStockMap, err
}
