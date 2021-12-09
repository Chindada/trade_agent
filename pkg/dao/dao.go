// Package dao package dao
package dao

import (
	"sync"
	"time"
	"trade_agent/pkg/config"
	"trade_agent/pkg/log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	gormlogger "gorm.io/gorm/logger"
)

var (
	globalConnection *gorm.DB
	once             sync.Once
)

func initConnection() {
	if globalConnection != nil {
		return
	}
	// logger for gorm
	dbLogger := gormlogger.New(log.Get(),
		gormlogger.Config{
			SlowThreshold:             1000 * time.Millisecond,
			Colorful:                  false,
			IgnoreRecordNotFoundError: false,
			LogLevel:                  gormlogger.Warn,
		})
	var err error
	var conf config.Config
	if conf, err = config.Get(); err != nil {
		log.Get().Panic(err)
	}
	dbSettings := conf.GetDBConfig()
	dsn := "host=" + dbSettings.DBHost + " user=" + dbSettings.DBUser +
		" password=" + dbSettings.DBPass + " dbname=" + dbSettings.Database +
		" port=" + dbSettings.DBPort + " sslmode=disable" +
		" TimeZone=" + dbSettings.DBTimeZone

	globalConnection, err = gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: dbLogger, SkipDefaultTransaction: true})
	if err != nil {
		log.Get().Panic(err)
	}
	// err = globalConnection.AutoMigrate(
	// 	&balance.Balance{},
	// 	&analyzeentiretick.AnalyzeEntireTick{},
	// 	&bidask.BidAsk{},
	// 	&entiretick.EntireTick{},
	// 	&holiday.Holiday{},
	// 	&kbar.Kbar{},
	// 	&simulate.Result{},
	// 	&simulationcond.AnalyzeCondition{},
	// 	&stock.Stock{},
	// 	&streamtick.StreamTick{},
	// 	&targetstock.Target{},
	// 	&tradeevent.EventResponse{},
	// 	&traderecord.TradeRecord{},
	// )
	if err != nil {
		log.Get().Panic(err)
	}
}

// Get Get
func Get() *gorm.DB {
	if globalConnection != nil {
		return globalConnection
	}
	once.Do(initConnection)
	return globalConnection
}
