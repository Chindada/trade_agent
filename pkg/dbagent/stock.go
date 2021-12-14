// Package dbagent package dbagent
package dbagent

import "gorm.io/gorm"

// Stock Stock
type Stock struct {
	gorm.Model `json:"-" swaggerignore:"true"`

	Number             string  `json:"number" gorm:"column:number;uniqueIndex"`
	Name               string  `json:"name" gorm:"column:name"`
	Exchange           string  `json:"exchange" gorm:"column:exchange"`
	Category           string  `json:"category" gorm:"column:category"`
	DayTrade           bool    `json:"day_trade" gorm:"column:day_trade"`
	LastClose          float64 `json:"last_close" gorm:"column:last_close"`
	LastVolume         int64   `json:"last_volume" gorm:"column:last_volume"`
	LastCloseChangePct float64 `json:"last_close_change_pct" gorm:"column:last_close_change_pct"`
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
func (c *DBAgent) InsertMultiStock(stockArr []*Stock) error {
	err := c.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.CreateInBatches(&stockArr, multiInsertBatchSize).Error; err != nil {
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
	result := c.DB.Model(&Stock{}).Not("name = ?", "").Order("number").Find(&stockArr)
	for _, v := range stockArr {
		allStockMap[v.Number] = v
	}
	return allStockMap, result.Error
}
