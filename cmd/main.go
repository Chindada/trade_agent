package main

import (
	"trade_agent/pkg/dbagent"
	"trade_agent/pkg/modules/stock"
	"trade_agent/pkg/modules/tradeevent"
	"trade_agent/pkg/mqhandler"
	"trade_agent/pkg/routers"
	"trade_agent/pkg/sinopacapi"
	"trade_agent/pkg/tasks"
)

func main() {
	// init don't move sequence
	dbagent.InitDatabase()
	tasks.InitTasks()
	sinopacapi.InitSinpacAPI()
	routers.ServeHTTP()

	// process chan
	keep := make(chan struct{})
	handler := mqhandler.Get()

	stock.InitStock(handler)
	tradeevent.InitTradeEvent(handler)

	<-keep
}
