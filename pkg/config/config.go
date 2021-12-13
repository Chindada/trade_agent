// Package config package config
package config

import (
	"io/ioutil"
	"sync"
	"trade_agent/pkg/log"

	"gopkg.in/yaml.v2"
)

var (
	globalConfig *Config
	once         sync.Once
)

// Config Config
type Config struct {
	Server   `yaml:"server"`
	Database `yaml:"database"`
	Schedule `yaml:"schedule"`
	MQTT     `yaml:"mqtt"`
	Trade    `yaml:"trade"`
}

// Server Server
type Server struct {
	RunMode        string `yaml:"run_mode"`
	HTTPPort       string `yaml:"http_port"`
	SinopacSRVHost string `yaml:"sinopac_srv_host"`
	SinopacSRVPort string `yaml:"sinopac_srv_port"`
}

// Database Database
type Database struct {
	DBHost     string `yaml:"db_host"`
	DBPort     string `yaml:"db_port"`
	DBUser     string `yaml:"db_user"`
	DBPass     string `yaml:"db_pass"`
	Database   string `yaml:"database"`
	DBEncode   string `yaml:"db_encode"`
	DBTimeZone string `yaml:"db_timezone"`
}

// MQTT MQTT
type MQTT struct {
	Host     string `yaml:"host" json:"host"`
	Port     string `yaml:"port" json:"port"`
	User     string `yaml:"user" json:"user"`
	Password string `yaml:"password" json:"password"`
	ClientID string `yaml:"client_id" json:"client_id"`
	CAPath   string `yaml:"ca_path" json:"ca_path"`
	CertPath string `yaml:"cert_path" json:"cert_path"`
	KeyPath  string `yaml:"key_path" json:"key_path"`
}

// Trade Trade
type Trade struct {
	KbarPeriod      string `yaml:"kbar_period"`
	TargetCondition string `yaml:"target_condition"`
	BlackStock      string `yaml:"black_stock"`
	BlackCategory   string `yaml:"black_category"`
}

// Schedule Schedule
type Schedule struct {
	CleaneventCron     string `yaml:"cleanevent_cron"`
	RestartSinopacCron string `yaml:"restart_sinopac_cron"`
}

// parseConfig parseConfig
func parseConfig() {
	if globalConfig != nil {
		return
	}
	yamlFile, err := ioutil.ReadFile("./configs/config.yaml")
	if err != nil {
		log.Get().Panic(err)
	}
	err = yaml.Unmarshal(yamlFile, &globalConfig)
	if err != nil {
		log.Get().Panic(err)
	}
}

// Get Get
func Get() (config Config, err error) {
	if globalConfig != nil {
		return *globalConfig, err
	}
	once.Do(parseConfig)
	return *globalConfig, err
}

// GetServerConfig GetServerConfig
func (c Config) GetServerConfig() Server {
	return c.Server
}

// GetDBConfig GetDBConfig
func (c Config) GetDBConfig() Database {
	return c.Database
}

// GetScheduleConfig GetScheduleConfig
func (c Config) GetScheduleConfig() Schedule {
	return c.Schedule
}

// GetTradeConfig GetTradeConfig
func (c Config) GetTradeConfig() Trade {
	return c.Trade
}

// GetMQConfig GetMQConfig
func (c Config) GetMQConfig() MQTT {
	return c.MQTT
}
