// Package config package config
package config

import "trade_agent/pkg/log"

// TurnTradeInSwitchOFF TurnTradeInSwitchOFF
func TurnTradeInSwitchOFF() {
	if globalConfig == nil {
		log.Get().Panic("config setter should after init")
	}
	globalConfig.Switch.Buy = false
	globalConfig.Switch.SellFirst = false
}
