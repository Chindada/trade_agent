package main

import (
	"trade_agent/pkg/dbagent"
	"trade_agent/pkg/modules/future"
	"trade_agent/pkg/modules/history"
	"trade_agent/pkg/modules/order"
	"trade_agent/pkg/modules/simulation"
	"trade_agent/pkg/modules/stock"
	"trade_agent/pkg/modules/subscribe"
	"trade_agent/pkg/modules/targets"
	"trade_agent/pkg/modules/tickprocess"
	"trade_agent/pkg/modules/tradeday"
	"trade_agent/pkg/modules/tradeevent"
	"trade_agent/pkg/mqhandler"
	"trade_agent/pkg/routers"
	"trade_agent/pkg/sinopacapi"
	"trade_agent/pkg/tasks"
)

func main() {
	// initial basic
	dbagent.InitDatabase()
	mqhandler.InitMQHandler()
	sinopacapi.InitSinpacAPI()
	tasks.InitTasks()

	// serve http before trade agent module start
	routers.ServeHTTP()
	// stuck chan
	keep := make(chan struct{})

	// trade agent modules
	order.InitOrder()
	tradeevent.InitTradeEvent()
	tickprocess.InitTickProcess()

	// should before targets
	subscribe.InitSubscribe()
	history.InitHistory()

	tradeday.InitTradeDay()
	stock.InitStock()
	future.InitFuture()
	targets.InitTargets()

	simulation.InitSimulation()

	// stuck
	<-keep
}
