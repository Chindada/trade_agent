// Package tasks package tasks
package tasks

import (
	"trade_agent/pkg/config"
	"trade_agent/pkg/log"
	"trade_agent/pkg/tasks/healthchecktask"
	"trade_agent/pkg/tasks/tradeeventtask"

	"github.com/robfig/cron"
)

// InitTasks InitTasks
func InitTasks() {
	log.Get().Info("Initial Tasks")
	var err error
	var conf config.Config

	conf, err = config.Get()
	if err != nil {
		log.Get().Panic(err)
	}
	schedule := conf.GetScheduleConfig()

	c := cron.New()
	err = c.AddFunc(schedule.CleaneventCron, func() {
		tradeeventtask.Run()
	})
	if err != nil {
		log.Get().Panic(err)
	}

	err = c.AddFunc(schedule.RestartSinopacCron, func() {
		healthchecktask.Run()
	})
	if err != nil {
		log.Get().Panic(err)
	}

	c.Start()
}
