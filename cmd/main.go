package main

import (
	"os"
	"trade_agent/global"
	"trade_agent/pkg/dbagent"
	"trade_agent/pkg/modules/cloudevent"
	"trade_agent/pkg/modules/history"
	"trade_agent/pkg/modules/order"
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
	// check if env is production or development
	deployment := os.Getenv("DEPLOYMENT")
	if deployment != "docker" {
		global.Development = true
	}

	keep := make(chan struct{})
	// initial core funcs
	dbagent.InitDatabase()
	mqhandler.InitMQHandler()
	sinopacapi.InitSinpacAPI()

	// init cron tasks
	tasks.InitTasks()

	// serve http before trade agent module start
	routers.ServeHTTP()

	// fill all basic data first
	tradeday.InitTradeDay()
	stock.InitStock()

	// sub sino srv event to update db event
	cloudevent.InitCloudEvent()

	// wait order and simulation result to place order
	// update all order status looply
	order.InitOrder()

	// process all tick
	// include realtime, history tick, kbar
	// wait simulation result to process tick
	tickprocess.InitTickProcess()

	// receive target to subscribe tick, bidask
	// receive target to fill history data
	subscribe.InitSubscribe()
	history.InitHistory()

	// find target
	targets.InitTargets()
	<-keep
}
