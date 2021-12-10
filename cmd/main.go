package main

import (
	"trade_agent/pkg/dao"
	"trade_agent/pkg/routers"
	"trade_agent/pkg/tasks"
)

func main() {
	dao.InitDatabase()
	tasks.InitTasks()
	routers.ServeHTTP()
	keep := make(chan struct{})
	<-keep
}
