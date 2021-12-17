// Package dbagent package dbagent
package dbagent

import "gorm.io/gorm"

// SimulationCondition SimulationCondition
type SimulationCondition struct {
	gorm.Model `json:"-" swaggerignore:"true"`

	TrimHistoryCloseCount bool    `json:"trim_history_close_count,omitempty" yaml:"trim_history_close_count" gorm:"column:trim_history_close_count"`
	HistoryCloseCount     int64   `json:"history_close_count,omitempty" yaml:"history_close_count" gorm:"column:history_close_count"`
	ForwardOutInRatio     float64 `json:"forward_out_in_ratio,omitempty" yaml:"forward_out_in_ratio" gorm:"column:forward_out_in_ratio"`
	ReverseOutInRatio     float64 `json:"reverse_out_in_ratio,omitempty" yaml:"reverse_out_in_ratio" gorm:"column:reverse_out_in_ratio"`
	CloseChangeRatioLow   float64 `json:"close_change_ratio_low,omitempty" yaml:"close_change_ratio_low" gorm:"column:close_change_ratio_low"`
	CloseChangeRatioHigh  float64 `json:"close_change_ratio_high,omitempty" yaml:"close_change_ratio_high" gorm:"column:close_change_ratio_high"`
	OpenChangeRatio       float64 `json:"open_change_ratio,omitempty" yaml:"open_change_ratio" gorm:"column:open_change_ratio"`
	RsiHigh               float64 `json:"rsi_high,omitempty" yaml:"rsi_high" gorm:"column:rsi_high"`
	RsiLow                float64 `json:"rsi_low,omitempty" yaml:"rsi_low" gorm:"column:rsi_low"`
	TicksPeriodThreshold  float64 `json:"ticks_period_threshold,omitempty" yaml:"ticks_period_threshold" gorm:"column:ticks_period_threshold"`
	TicksPeriodLimit      float64 `json:"ticks_period_limit,omitempty" yaml:"ticks_period_limit" gorm:"column:ticks_period_limit"`
	TicksPeriodCount      int     `json:"ticks_period_count,omitempty" yaml:"ticks_period_count" gorm:"column:ticks_period_count"`
	VolumePerSecond       int64   `json:"volume_per_second,omitempty" yaml:"volume_per_second" gorm:"column:volume_per_second"`
}

// InsertSimulationCondition InsertSimulationCondition
func (c *DBAgent) InsertSimulationCondition(record *SimulationCondition) error {
	err := c.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&record).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// InsertMultiSimulationCondition InsertMultiSimulationCondition
func (c *DBAgent) InsertMultiSimulationCondition(recordArr []*SimulationCondition) error {
	err := c.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.CreateInBatches(&recordArr, multiInsertBatchSize).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// DeleteAllSimulationCondition DeleteAllSimulationCondition
func (c *DBAgent) DeleteAllSimulationCondition() error {
	err := c.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Not("id = 0").Unscoped().Delete(&SimulationCondition{}).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}
