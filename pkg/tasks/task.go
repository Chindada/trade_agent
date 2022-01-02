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
	log.Get().Info("Initial Tasks")

	var err error
	schedule := config.GetScheduleConfig()

	c := cron.New()
	err = c.AddFunc(schedule.CleanEvent, func() {
		cloudeventtask.Run()
	})
	if err != nil {
		log.Get().Panic(err)
	}

	err = c.AddFunc(schedule.RestartSinopac, func() {
		healthchecktask.Run()
	})
	if err != nil {
		log.Get().Panic(err)
	}

	c.Start()
}
