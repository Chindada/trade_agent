// Package dbagent package dbagent
package dbagent

import (
	"time"

	"gorm.io/gorm"
)

// RealTimeBidAsk RealTimeBidAsk
type RealTimeBidAsk struct {
	gorm.Model `json:"-" swaggerignore:"true"`

	Stock    *Stock    `json:"stock,omitempty" yaml:"stock" gorm:"foreignKey:StockID"`
	StockID  int64     `json:"stock_id,omitempty" yaml:"stock_id" gorm:"column:stock_id"`
	TickTime time.Time `json:"tick_time,omitempty" yaml:"tick_time" gorm:"column:tick_time"`

	BidPrice1   float64 `json:"bid_price_1,omitempty" yaml:"bid_price_1" gorm:"column:bid_price_1"`
	BidVolume1  int64   `json:"bid_volume_1,omitempty" yaml:"bid_volume_1" gorm:"column:bid_volume_1"`
	DiffBidVol1 int64   `json:"diff_bid_vol_1,omitempty" yaml:"diff_bid_vol_1" gorm:"column:diff_bid_vol_1"`
	BidPrice2   float64 `json:"bid_price_2,omitempty" yaml:"bid_price_2" gorm:"column:bid_price_2"`
	BidVolume2  int64   `json:"bid_volume_2,omitempty" yaml:"bid_volume_2" gorm:"column:bid_volume_2"`
	DiffBidVol2 int64   `json:"diff_bid_vol_2,omitempty" yaml:"diff_bid_vol_2" gorm:"column:diff_bid_vol_2"`
	BidPrice3   float64 `json:"bid_price_3,omitempty" yaml:"bid_price_3" gorm:"column:bid_price_3"`
	BidVolume3  int64   `json:"bid_volume_3,omitempty" yaml:"bid_volume_3" gorm:"column:bid_volume_3"`
	DiffBidVol3 int64   `json:"diff_bid_vol_3,omitempty" yaml:"diff_bid_vol_3" gorm:"column:diff_bid_vol_3"`
	BidPrice4   float64 `json:"bid_price_4,omitempty" yaml:"bid_price_4" gorm:"column:bid_price_4"`
	BidVolume4  int64   `json:"bid_volume_4,omitempty" yaml:"bid_volume_4" gorm:"column:bid_volume_4"`
	DiffBidVol4 int64   `json:"diff_bid_vol_4,omitempty" yaml:"diff_bid_vol_4" gorm:"column:diff_bid_vol_4"`
	BidPrice5   float64 `json:"bid_price_5,omitempty" yaml:"bid_price_5" gorm:"column:bid_price_5"`
	BidVolume5  int64   `json:"bid_volume_5,omitempty" yaml:"bid_volume_5" gorm:"column:bid_volume_5"`
	DiffBidVol5 int64   `json:"diff_bid_vol_5,omitempty" yaml:"diff_bid_vol_5" gorm:"column:diff_bid_vol_5"`
	AskPrice1   float64 `json:"ask_price_1,omitempty" yaml:"ask_price_1" gorm:"column:ask_price_1"`
	AskVolume1  int64   `json:"ask_volume_1,omitempty" yaml:"ask_volume_1" gorm:"column:ask_volume_1"`
	DiffAskVol1 int64   `json:"diff_ask_vol_1,omitempty" yaml:"diff_ask_vol_1" gorm:"column:diff_ask_vol_1"`
	AskPrice2   float64 `json:"ask_price_2,omitempty" yaml:"ask_price_2" gorm:"column:ask_price_2"`
	AskVolume2  int64   `json:"ask_volume_2,omitempty" yaml:"ask_volume_2" gorm:"column:ask_volume_2"`
	DiffAskVol2 int64   `json:"diff_ask_vol_2,omitempty" yaml:"diff_ask_vol_2" gorm:"column:diff_ask_vol_2"`
	AskPrice3   float64 `json:"ask_price_3,omitempty" yaml:"ask_price_3" gorm:"column:ask_price_3"`
	AskVolume3  int64   `json:"ask_volume_3,omitempty" yaml:"ask_volume_3" gorm:"column:ask_volume_3"`
	DiffAskVol3 int64   `json:"diff_ask_vol_3,omitempty" yaml:"diff_ask_vol_3" gorm:"column:diff_ask_vol_3"`
	AskPrice4   float64 `json:"ask_price_4,omitempty" yaml:"ask_price_4" gorm:"column:ask_price_4"`
	AskVolume4  int64   `json:"ask_volume_4,omitempty" yaml:"ask_volume_4" gorm:"column:ask_volume_4"`
	DiffAskVol4 int64   `json:"diff_ask_vol_4,omitempty" yaml:"diff_ask_vol_4" gorm:"column:diff_ask_vol_4"`
	AskPrice5   float64 `json:"ask_price_5,omitempty" yaml:"ask_price_5" gorm:"column:ask_price_5"`
	AskVolume5  int64   `json:"ask_volume_5,omitempty" yaml:"ask_volume_5" gorm:"column:ask_volume_5"`
	DiffAskVol5 int64   `json:"diff_ask_vol_5,omitempty" yaml:"diff_ask_vol_5" gorm:"column:diff_ask_vol_5"`

	Suspend  int64 `json:"suspend,omitempty" yaml:"suspend" gorm:"column:suspend"`
	Simtrade int64 `json:"simtrade,omitempty" yaml:"simtrade" gorm:"column:simtrade"`
}

// TableName TableName
func (RealTimeBidAsk) TableName() string {
	return "realtime_bidask"
}

// InsertRealTimeBidAsk InsertRealTimeBidAsk
func (c *DBAgent) InsertRealTimeBidAsk(record *RealTimeBidAsk) error {
	err := c.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&record).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// InsertMultiRealTimeBidAsk InsertMultiRealTimeBidAsk
func (c *DBAgent) InsertMultiRealTimeBidAsk(recordArr []*RealTimeBidAsk) error {
	err := c.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.CreateInBatches(&recordArr, multiInsertBatchSize).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// DeleteAllRealTimeBidAsk DeleteAllRealTimeBidAsk
func (c *DBAgent) DeleteAllRealTimeBidAsk() error {
	err := c.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Not("id = 0").Unscoped().Delete(&RealTimeBidAsk{}).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}
