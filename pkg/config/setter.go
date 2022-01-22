// Package config package config
package config

import "trade_agent/pkg/log"

// TurnTradeInSwitchOFF TurnTradeInSwitchOFF
func TurnTradeInSwitchOFF() {
	if globalConfig == nil {
		log.Get().Panic("config setter should after init")
	}
	globalConfig.TradeSwitch.Buy = false
	globalConfig.TradeSwitch.SellFirst = false
}
