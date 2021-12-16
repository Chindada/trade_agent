// Package dbagent package dbagent
package dbagent

import (
	"time"

	"gorm.io/gorm"
)

// SimulationResult SimulationResult
type SimulationResult struct {
	gorm.Model `json:"-" swaggerignore:"true"`

	Balance        int64     `json:"balance,omitempty" yaml:"balance" gorm:"column:balance"`
	ForwardBalance int64     `json:"forward_balance,omitempty" yaml:"forward_balance" gorm:"column:forward_balance"`
	ReverseBalance int64     `json:"reverse_balance,omitempty" yaml:"reverse_balance" gorm:"column:reverse_balance"`
	TotalLoss      int64     `json:"total_loss,omitempty" yaml:"total_loss" gorm:"column:total_loss"`
	TradeCount     int64     `json:"trade_count,omitempty" yaml:"trade_count" gorm:"column:trade_count"`
	PositiveDays   int64     `json:"positive_days,omitempty" yaml:"positive_days" gorm:"column:positive_days"`
	NegativeDays   int64     `json:"negative_days,omitempty" yaml:"negative_days" gorm:"column:negative_days"`
	TotalDays      int64     `json:"total_days,omitempty" yaml:"total_days" gorm:"column:total_days"`
	IsBestForward  bool      `json:"is_best_forward,omitempty" yaml:"is_best_forward" gorm:"column:is_best_forward"`
	IsBestReverse  bool      `json:"is_best_reverse,omitempty" yaml:"is_best_reverse" gorm:"column:is_best_reverse"`
	TradeDay       time.Time `json:"trade_day,omitempty" yaml:"trade_day" gorm:"column:trade_day"`

	Cond   *SimulationCondition `json:"cond,omitempty" yaml:"cond" gorm:"foreignKey:CondID"`
	CondID int64                `json:"cond_id,omitempty" yaml:"cond_id" gorm:"column:cond_id"`
}
