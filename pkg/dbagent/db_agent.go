// Package dbagent package dbagent
package dbagent

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	"trade_agent/pkg/config"
	"trade_agent/pkg/log"

	_ "github.com/lib/pq" // postgres driver for "database/sql"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	gormlogger "gorm.io/gorm/logger"
)

// Tabler Tabler
type Tabler interface {
	TableName() string
}

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
	log.Get().Info("Initial Database")

	dbSettings := config.GetDatabaseConfig()
	dsn := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s sslmode=disable TimeZone=%s",
		dbSettings.User,
		dbSettings.Passwd,
		dbSettings.Host,
		dbSettings.Port,
		dbSettings.TimeZone,
	)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Get().Panic(err)
	}
	defer func() {
		if err = db.Close(); err != nil {
			log.Get().Panic(err)
		}
	}()
	var exist bool
	err = db.QueryRow(fmt.Sprintf("SELECT EXISTS(SELECT datname FROM pg_catalog.pg_database WHERE datname = '%s')", dbSettings.Database)).Scan(&exist)
	if err != nil {
		log.Get().Panic(err)
	}
	if !exist {
		_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbSettings.Database))
		if err != nil {
			log.Get().Panic(err)
		}
	}
	once.Do(initConnection)
}

func initConnection() {
	if globalAgent != nil {
		return
	}
	// logger for gorm
	dbLogger := gormlogger.New(log.Get(),
		gormlogger.Config{
			SlowThreshold:             1500 * time.Millisecond,
			Colorful:                  false,
			IgnoreRecordNotFoundError: true,
			LogLevel:                  gormlogger.Warn,
		})
	var err error
	dbSettings := config.GetDatabaseConfig()
	dsn := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=disable TimeZone=%s",
		dbSettings.User,
		dbSettings.Passwd,
		dbSettings.Database,
		dbSettings.Host,
		dbSettings.Port,
		dbSettings.TimeZone,
	)

	var newAgent DBAgent
	newAgent.DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: dbLogger, SkipDefaultTransaction: true})
	if err != nil {
		log.Get().Panic(err)
	}
	err = newAgent.DB.AutoMigrate(
		&Balance{},
		&CalendarDate{},
		&CloudEvent{},
		&HistoryClose{},
		&HistoryKbar{},
		&HistoryTick{},
		&OrderStatus{},
		&RealTimeBidAsk{},
		&RealTimeTick{},
		&Stock{},
		&Target{},
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
