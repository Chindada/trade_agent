// Package dbagent package dbagent
package dbagent

import "gorm.io/gorm"

// Stock Stock
type Stock struct {
	gorm.Model `json:"-" swaggerignore:"true"`

	Number    string  `json:"number,omitempty" yaml:"number" gorm:"column:number"`
	Name      string  `json:"name,omitempty" yaml:"name" gorm:"column:name"`
	Exchange  string  `json:"exchange,omitempty" yaml:"exchange" gorm:"column:exchange"`
	Category  string  `json:"category,omitempty" yaml:"category" gorm:"column:category"`
	DayTrade  bool    `json:"day_trade,omitempty" yaml:"day_trade" gorm:"column:day_trade"`
	LastClose float64 `json:"last_close,omitempty" yaml:"last_close" gorm:"column:last_close"`
}

// TableName TableName
func (Stock) TableName() string {
	return "basic_stock"
}

// InsertStock InsertStock
func (c *DBAgent) InsertStock(record *Stock) error {
	err := c.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&record).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// InsertMultiStock InsertMultiStock
func (c *DBAgent) InsertMultiStock(recordArr []*Stock) error {
	err := c.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.CreateInBatches(&recordArr, multiInsertBatchSize).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// InsertOrUpdateMultiStock InsertOrUpdateMultiStock
func (c *DBAgent) InsertOrUpdateMultiStock(recordArr []*Stock) error {
	err := c.DB.Transaction(func(tx *gorm.DB) error {
		for _, stock := range recordArr {
			tmp := stock
			if tmp.ID != 0 {
				if err := tx.Save(&tmp).Error; err != nil {
					return err
				}
			} else {
				if err := tx.Create(&tmp).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
	return err
}

// UpdateAllStockDayTradeFalse UpdateAllStockDayTradeFalse
func (c *DBAgent) UpdateAllStockDayTradeFalse() error {
	err := c.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&Stock{}).Where("id != 0").Updates(Stock{DayTrade: false}).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// DeleteAllStock DeleteAllStock
func (c *DBAgent) DeleteAllStock() error {
	err := c.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Not("id = 0").Unscoped().Delete(&Stock{}).Error; err != nil {
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
	err = c.DB.Model(&Stock{}).Not("name = ?", "").Not("category = ?", "00").Not("category = ?", "custom").Order("number").Find(&stockArr).Error

	for _, v := range stockArr {
		allStockMap[v.Number] = v
	}
	return allStockMap, err
}

// GetAllDayTradeStockMap GetAllDayTradeStockMap
func (c *DBAgent) GetAllDayTradeStockMap() (allStockMap map[string]*Stock, err error) {
	allStockMap = make(map[string]*Stock)
	var stockArr []*Stock
	err = c.DB.Model(&Stock{}).Where("day_trade = ?", true).Find(&stockArr).Error

	for _, v := range stockArr {
		allStockMap[v.Number] = v
	}
	return allStockMap, err
}

// GetAllCustomStockMap GetAllCustomStockMap
func (c *DBAgent) GetAllCustomStockMap() (customMap map[string]*Stock, err error) {
	customMap = make(map[string]*Stock)
	var stockArr []*Stock
	err = c.DB.Model(&Stock{}).Where("category = ?", "custom").Find(&stockArr).Error

	for _, v := range stockArr {
		customMap[v.Number] = v
	}
	return customMap, err
}
