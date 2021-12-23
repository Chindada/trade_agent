// Package dbagent package dbagent
package dbagent

import (
	"database/sql"
	"sync"
	"time"
	"trade_agent/pkg/config"
	"trade_agent/pkg/log"

	_ "github.com/lib/pq" // postgres driver for "database/sql"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	gormlogger "gorm.io/gorm/logger"
)

// DBAgent DBAgent
type DBAgent struct {
	DB *gorm.DB
}

const (
	multiInsertBatchSize int = 2000
)

var (
	globalAgent *DBAgent
	once        sync.Once
)

// InitDatabase InitDatabase
func InitDatabase() {
	conf, err := config.Get()
	if err != nil {
		log.Get().Panic(err)
	}
	dbSettings := conf.GetDBConfig()
	db, err := sql.Open(
		"postgres",
		"user="+dbSettings.DBUser+" password="+dbSettings.DBPass+" host="+dbSettings.DBHost+" port="+dbSettings.DBPort+" sslmode=disable")
	if err != nil {
		log.Get().Panic(err)
	}
	defer func() {
		if err = db.Close(); err != nil {
			log.Get().Panic(err)
		}
	}()
	var exist bool
	statement := "SELECT EXISTS(SELECT datname FROM pg_catalog.pg_database WHERE datname = '" + dbSettings.Database + "')"
	err = db.QueryRow(statement).Scan(&exist)
	if err != nil {
		log.Get().Panic(err)
	}
	if !exist {
		_, err = db.Exec("CREATE DATABASE " + dbSettings.Database + ";")
		if err != nil {
			log.Get().Panic(err)
		}
	}
	once.Do(initConnection)
	log.Get().Info("Initial Database")
}

func initConnection() {
	if globalAgent != nil {
		return
	}
	// logger for gorm
	dbLogger := gormlogger.New(log.Get(),
		gormlogger.Config{
			SlowThreshold:             1000 * time.Millisecond,
			Colorful:                  false,
			IgnoreRecordNotFoundError: true,
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

	var newAgent DBAgent
	newAgent.DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: dbLogger, SkipDefaultTransaction: true})
	if err != nil {
		log.Get().Panic(err)
	}
	err = newAgent.DB.AutoMigrate(
		&AnalyzedTick{},
		&Balance{},
		&CalendarDate{},
		&RealTimeBidAsk{},
		&HistoryKbar{},
		&HistoryTick{},
		&HistoryClose{},
		&RealTimeTick{},
		&Stock{},
		&Target{},
		&CloudEvent{},
		&OrderStatus{},
		&SimulationResult{},
		&SimulationCondition{},
	)
	if err != nil {
		log.Get().Panic(err)
	}
	globalAgent = &newAgent
}

// Get Get
func Get() *DBAgent {
	if globalAgent != nil {
		return globalAgent
	}
	log.Get().Panic("globalAgent was not initailized")
	return nil
}
