package main

import (
	"trade_agent/pkg/dbagent"
	"trade_agent/pkg/modules/cloudevent"
	"trade_agent/pkg/modules/future"
	"trade_agent/pkg/modules/history"
	"trade_agent/pkg/modules/order"
	"trade_agent/pkg/modules/simulation"
	"trade_agent/pkg/modules/stock"
	"trade_agent/pkg/modules/subscribe"
	"trade_agent/pkg/modules/targets"
	"trade_agent/pkg/modules/tickprocess"
	"trade_agent/pkg/modules/tradeday"
	"trade_agent/pkg/mqhandler"
	"trade_agent/pkg/routers"
	"trade_agent/pkg/sinopacapi"
	"trade_agent/pkg/tasks"
)

func main() {
	keep := make(chan struct{})

	// initial core utils
	dbagent.InitDatabase()
	mqhandler.InitMQHandler()
	sinopacapi.InitSinpacAPI()

	// init cron tasks
	tasks.InitTasks()

	// serve http before trade agent module start
	routers.ServeHTTP()

	// fill all basic data
	tradeday.InitTradeDay()
	stock.InitStock()
	future.InitFuture()

	// wait order and simulation result to place order
	// update all order status looply
	order.InitOrder()

	// simulation
	simulation.InitSimulation()

	// update sino srv event
	cloudevent.InitCloudEvent()

	// process all tick
	// include realtime, history tick, kbar
	tickprocess.InitTickProcess()

	// receive target to subscribe tick, bidask
	subscribe.InitSubscribe()

	// receive target to fill history data
	history.InitHistory()

	// find target
	targets.InitTargets()

	<-keep
}
