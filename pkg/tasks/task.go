// Package tasks package tasks
package tasks

import (
	"trade_agent/pkg/config"
	"trade_agent/pkg/log"
	"trade_agent/pkg/tasks/cloudeventtask"
	"trade_agent/pkg/tasks/healthchecktask"

	"github.com/robfig/cron"
)

// InitTasks InitTasks
func InitTasks() {
	var err error
	var conf config.Config

	conf, err = config.Get()
	if err != nil {
		log.Get().Panic(err)
	}
	schedule := conf.GetScheduleConfig()

	c := cron.New()
	err = c.AddFunc(schedule.CleaneventCron, func() {
		cloudeventtask.Run()
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
	log.Get().Info("Initial Tasks")
}
