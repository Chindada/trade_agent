// Package dao package dao
package dao

import (
	"database/sql"
	"trade_agent/pkg/config"
	"trade_agent/pkg/log"

	_ "github.com/lib/pq" // postgres driver for "database/sql"
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
}
