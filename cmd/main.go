package main

import (
	"trade_agent/pkg/dbagent"
	"trade_agent/pkg/modules/analyze"
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

	// trade agent modules
	keep := make(chan struct{})

	// start modules
	tradeevent.InitTradeEvent()
	order.InitOrder()
	tickprocess.InitTickProcess()
	subscribe.InitSubscribe()
	tradeday.InitTradeDay()
	history.InitHistory()
	stock.InitStock()
	targets.InitTargets()

	future.InitFuture()
	analyze.InitAnalyze()
	simulation.InitSimulation()

	// stuck
	<-keep
}
