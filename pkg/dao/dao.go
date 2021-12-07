// Package dao package dao
package dao

import (
	"sync"
	"time"

	"gitlab.tocraw.com/root/toc_trader/pkg/config"
	"gitlab.tocraw.com/root/toc_trader/pkg/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	gormlogger "gorm.io/gorm/logger"
)

var (
	globalConnection *gorm.DB
	initLock         sync.Mutex
)

func initConnection() {
	defer initLock.Unlock()
	initLock.Lock()
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
	var allConfigs config.Config
	if allConfigs, err = config.Get(); err != nil {
		log.Get().Panic(err)
	}

	dbSettings := allConfigs.GetDBConfig()
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
	initConnection()
	return globalConnection
}
