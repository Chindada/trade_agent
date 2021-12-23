// Package tickprocess package tickprocess
package tickprocess

import "trade_agent/pkg/log"

// InitTickProcess InitTickProcess
func InitTickProcess() {
	err := subRealTimeTick()
	if err != nil {
		log.Get().Panic(err)
	}
	err = subRealTimeBidAsk()
	if err != nil {
		log.Get().Panic(err)
	}
	err = subHistroyTick()
	if err != nil {
		log.Get().Panic(err)
	}
	err = subHistroyKbar()
	if err != nil {
		log.Get().Panic(err)
	}
	log.Get().Info("Initial TickProcess")
}
