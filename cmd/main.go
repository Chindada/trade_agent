package main

import (
	"os"
	"path/filepath"

	"trade_agent/global"
	"trade_agent/pkg/dbagent"
	"trade_agent/pkg/log"
	"trade_agent/pkg/modules/analyze"
	"trade_agent/pkg/modules/cloudevent"
	"trade_agent/pkg/modules/healthcheck"
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

func init() {
	// get binary path
	ex, err := os.Executable()
	if err != nil {
		log.Get().Panic(err)
	}
	global.Get().SetBasePath(filepath.Clean(filepath.Dir(ex)))

	// check if env is production or development
	deployment, exist := os.LookupEnv("DEPLOYMENT")
	if deployment != "docker" || !exist {
		global.Get().SetIsDevelopment(true)
	}
}

func main() {
	// initial core funcs
	dbagent.InitDatabase()
	mqhandler.InitMQHandler()
	sinopacapi.InitSinpacAPI()
	tasks.InitTasks()
	routers.ServeHTTP()

	// fill all basic module first
	tradeday.InitTradeDay()
	stock.InitStock()

	// sub sino srv event to update db event
	cloudevent.InitCloudEvent()

	// after history data is ready, sub to analyze
	analyze.InitAnalyze()

	// wait order and simulation result to place order
	// update all order status looply
	order.InitOrder()

	// receive target to subscribe tick, bidask
	subscribe.InitSubscribe()

	// process all tick
	// include realtime, history tick, kbar
	// wait simulation result to process tick
	tickprocess.InitTickProcess()

	// receive target to fill history data
	history.InitHistory()

	// find target
	targets.InitTargets()

	// initial health check and lock main here
	healthcheck.InitHealthCheck()
}
