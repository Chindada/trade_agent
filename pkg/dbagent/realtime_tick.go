// Package dbagent package dbagent
package dbagent

import (
	"time"

	"gorm.io/gorm"
)

// RealTimeTick RealTimeTick
type RealTimeTick struct {
	gorm.Model `json:"-" swaggerignore:"true"`

	Stock    *Stock    `json:"stock,omitempty" yaml:"stock" gorm:"foreignKey:StockID"`
	StockID  int64     `json:"stock_id,omitempty" yaml:"stock_id" gorm:"column:stock_id"`
	TickTime time.Time `json:"tick_time,omitempty" yaml:"tick_time" gorm:"column:tick_time"`

	Open            float64 `json:"open,omitempty" yaml:"open" gorm:"column:open"`
	AvgPrice        float64 `json:"avg_price,omitempty" yaml:"avg_price" gorm:"column:avg_price"`
	Close           float64 `json:"close,omitempty" yaml:"close" gorm:"column:close"`
	High            float64 `json:"high,omitempty" yaml:"high" gorm:"column:high"`
	Low             float64 `json:"low,omitempty" yaml:"low" gorm:"column:low"`
	Amount          float64 `json:"amount,omitempty" yaml:"amount" gorm:"column:amount"`
	AmountSum       float64 `json:"amount_sum,omitempty" yaml:"amount_sum" gorm:"column:amount_sum"`
	Volume          int64   `json:"volume,omitempty" yaml:"volume" gorm:"column:volume"`
	VolumeSum       int64   `json:"volume_sum,omitempty" yaml:"volume_sum" gorm:"column:volume_sum"`
	TickType        int64   `json:"tick_type,omitempty" yaml:"tick_type" gorm:"column:tick_type"`
	ChgType         int64   `json:"chg_type,omitempty" yaml:"chg_type" gorm:"column:chg_type"`
	PriceChg        float64 `json:"price_chg,omitempty" yaml:"price_chg" gorm:"column:price_chg"`
	PctChg          float64 `json:"pct_chg,omitempty" yaml:"pct_chg" gorm:"column:pct_chg"`
	BidSideTotalVol int64   `json:"bid_side_total_vol,omitempty" yaml:"bid_side_total_vol" gorm:"column:bid_side_total_vol"`
	AskSideTotalVol int64   `json:"ask_side_total_vol,omitempty" yaml:"ask_side_total_vol" gorm:"column:ask_side_total_vol"`
	BidSideTotalCnt int64   `json:"bid_side_total_cnt,omitempty" yaml:"bid_side_total_cnt" gorm:"column:bid_side_total_cnt"`
	AskSideTotalCnt int64   `json:"ask_side_total_cnt,omitempty" yaml:"ask_side_total_cnt" gorm:"column:ask_side_total_cnt"`

	Suspend  int64 `json:"suspend,omitempty" yaml:"suspend" gorm:"column:suspend"`
	Simtrade int64 `json:"simtrade,omitempty" yaml:"simtrade" gorm:"column:simtrade"`
}

// InsertRealTimeTick InsertRealTimeTick
func (c *DBAgent) InsertRealTimeTick(record *RealTimeTick) error {
	err := c.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&record).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// InsertMultiRealTimeTick InsertMultiRealTimeTick
func (c *DBAgent) InsertMultiRealTimeTick(recordArr []*RealTimeTick) error {
	err := c.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.CreateInBatches(&recordArr, multiInsertBatchSize).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// DeleteAllRealTimeTick DeleteAllRealTimeTick
func (c *DBAgent) DeleteAllRealTimeTick() error {
	err := c.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Not("id = 0").Unscoped().Delete(&RealTimeTick{}).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// RealTimeTickArr RealTimeTickArr
type RealTimeTickArr []*RealTimeTick

// GetStockNum GetStockNum
func (c RealTimeTickArr) GetStockNum() string {
	if len(c) != 0 {
		return c[0].Stock.Number
	}
	return ""
}

// GetTotalTime GetTotalTime
func (c RealTimeTickArr) GetTotalTime() float64 {
	if len(c) > 1 {
		return (float64(c[len(c)-1].TickTime.UnixNano()) - float64(c[0].TickTime.UnixNano())) / 1000 / 1000 / 1000
	}
	return 0
}

// Analyzer Analyzer
func (c RealTimeTickArr) Analyzer() []int64 {
	var analyzeTickArr RealTimeTickArr
	var volumeArr []int64
	for i, tick := range c {
		if i == 0 {
			continue
		}
		if len(analyzeTickArr) > 1 {
			if analyzeTickArr.GetTotalTime() > 5 {
				var volumeSum int64
				for _, k := range analyzeTickArr {
					volumeSum += k.Volume
				}
				analyzeTickArr = RealTimeTickArr{}
				volumeArr = append(volumeArr, volumeSum)
			}
		}
		analyzeTickArr = append(analyzeTickArr, tick)
	}
	return volumeArr
}
