// Package dbagent package dbagent
package dbagent

import (
	"time"

	"gorm.io/gorm"
)

// SimulationResult SimulationResult
type SimulationResult struct {
	gorm.Model `json:"-" swaggerignore:"true"`

	SimulationResult        int64     `json:"SimulationResult,omitempty" yaml:"SimulationResult" gorm:"column:SimulationResult"`
	ForwardSimulationResult int64     `json:"forward_SimulationResult,omitempty" yaml:"forward_SimulationResult" gorm:"column:forward_SimulationResult"`
	ReverseSimulationResult int64     `json:"reverse_SimulationResult,omitempty" yaml:"reverse_SimulationResult" gorm:"column:reverse_SimulationResult"`
	TotalLoss               int64     `json:"total_loss,omitempty" yaml:"total_loss" gorm:"column:total_loss"`
	TradeCount              int64     `json:"trade_count,omitempty" yaml:"trade_count" gorm:"column:trade_count"`
	PositiveDays            int64     `json:"positive_days,omitempty" yaml:"positive_days" gorm:"column:positive_days"`
	NegativeDays            int64     `json:"negative_days,omitempty" yaml:"negative_days" gorm:"column:negative_days"`
	TotalDays               int64     `json:"total_days,omitempty" yaml:"total_days" gorm:"column:total_days"`
	IsBestForward           bool      `json:"is_best_forward,omitempty" yaml:"is_best_forward" gorm:"column:is_best_forward"`
	IsBestReverse           bool      `json:"is_best_reverse,omitempty" yaml:"is_best_reverse" gorm:"column:is_best_reverse"`
	TradeDay                time.Time `json:"trade_day,omitempty" yaml:"trade_day" gorm:"column:trade_day"`

	Cond   *SimulationCondition `json:"cond,omitempty" yaml:"cond" gorm:"foreignKey:CondID"`
	CondID int64                `json:"cond_id,omitempty" yaml:"cond_id" gorm:"column:cond_id"`
}

// InsertSimulationResult InsertSimulationResult
func (c *DBAgent) InsertSimulationResult(record *SimulationResult) error {
	err := c.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&record).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// InsertMultiSimulationResult InsertMultiSimulationResult
func (c *DBAgent) InsertMultiSimulationResult(recordArr []*SimulationResult) error {
	err := c.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.CreateInBatches(&recordArr, multiInsertBatchSize).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// DeleteAllSimulationResult DeleteAllSimulationResult
func (c *DBAgent) DeleteAllSimulationResult() error {
	err := c.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Not("id = 0").Unscoped().Delete(&SimulationResult{}).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}
