// Package config package config
package config

import (
	"fmt"
	"io/ioutil"
	"sync"
	"trade_agent/global"
	"trade_agent/pkg/log"
	"trade_agent/pkg/utils"

	"gopkg.in/yaml.v2"
)

var (
	globalConfig *Config
	once         sync.Once
)

// Config Config
type Config struct {
	basePath   string
	Server     Server     `json:"server,omitempty" yaml:"server"`
	Database   Database   `json:"database,omitempty" yaml:"database"`
	Schedule   Schedule   `json:"schedule,omitempty" yaml:"schedule"`
	MQTT       MQTT       `json:"mqtt,omitempty" yaml:"mqtt"`
	Trade      Trade      `json:"trade,omitempty" yaml:"trade"`
	Switch     Switch     `json:"switch,omitempty" yaml:"switch"`
	TargetCond TargetCond `json:"target_cond,omitempty" yaml:"target_cond"`
}

// Server Server
type Server struct {
	RunMode        string `json:"run_mode,omitempty" yaml:"run_mode"`
	HTTPPort       string `json:"http_port,omitempty" yaml:"http_port"`
	SinopacSRVHost string `json:"sinopac_srv_host,omitempty" yaml:"sinopac_srv_host"`
	SinopacSRVPort string `json:"sinopac_srv_port,omitempty" yaml:"sinopac_srv_port"`
}

// Database Database
type Database struct {
	DBHost     string `json:"db_host,omitempty" yaml:"db_host"`
	DBPort     string `json:"db_port,omitempty" yaml:"db_port"`
	DBUser     string `json:"db_user,omitempty" yaml:"db_user"`
	DBPass     string `json:"db_pass,omitempty" yaml:"db_pass"`
	Database   string `json:"database,omitempty" yaml:"database"`
	DBEncode   string `json:"db_encode,omitempty" yaml:"db_encode"`
	DBTimeZone string `json:"db_time_zone,omitempty" yaml:"db_time_zone"`
}

// MQTT MQTT
type MQTT struct {
	Host     string `json:"host,omitempty" yaml:"host"`
	Port     string `json:"port,omitempty" yaml:"port"`
	User     string `json:"user,omitempty" yaml:"user"`
	Password string `json:"password,omitempty" yaml:"password"`
	ClientID string `json:"client_id,omitempty" yaml:"client_id"`
	CAPath   string `json:"ca_path,omitempty" yaml:"ca_path"`
	CertPath string `json:"cert_path,omitempty" yaml:"cert_path"`
	KeyPath  string `json:"key_path,omitempty" yaml:"key_path"`
}

// Trade Trade
type Trade struct {
	HistoryClosePeriod int64 `json:"history_close_period,omitempty" yaml:"history_close_period"`
	HistoryTickPeriod  int64 `json:"history_tick_period,omitempty" yaml:"history_tick_period"`
	HistoryKbarPeriod  int64 `json:"history_kbar_period,omitempty" yaml:"history_kbar_period"`
	TradeInWaitTime    int64 `json:"trade_in_wait_time,omitempty" yaml:"trade_in_wait_time"`
	TradeOutWaitTime   int64 `json:"trade_out_wait_time,omitempty" yaml:"trade_out_wait_time"`
	HoldMaxTime        int64 `json:"hold_max_time,omitempty" yaml:"hold_max_time"`
	WaitInOpen         int64 `json:"wait_in_open,omitempty" yaml:"wait_in_open"`
	TradeInEndTime     int64 `json:"trade_in_end_time,omitempty" yaml:"trade_in_end_time"`
}

// Switch Switch
type Switch struct {
	EnableBuy       bool  `json:"enable_buy,omitempty" yaml:"enable_buy"`
	EnableSell      bool  `json:"enable_sell,omitempty" yaml:"enable_sell"`
	EnableSellFirst bool  `json:"enable_sell_first,omitempty" yaml:"enable_sell_first"`
	EnableBuyLater  bool  `json:"enable_buy_later,omitempty" yaml:"enable_buy_later"`
	MeanTimeForward int64 `json:"mean_time_forward,omitempty" yaml:"mean_time_forward"`
	MeanTimeReverse int64 `json:"mean_time_reverse,omitempty" yaml:"mean_time_reverse"`
	ForwardMax      int64 `json:"forward_max,omitempty" yaml:"forward_max"`
	ReverseMax      int64 `json:"reverse_max,omitempty" yaml:"reverse_max"`
}

// TargetCond TargetCond
type TargetCond struct {
	LimitPriceLow  float64  `json:"limit_price_low,omitempty" yaml:"limit_price_low"`
	LimitPriceHigh float64  `json:"limit_price_high,omitempty" yaml:"limit_price_high"`
	LimitVolume    int64    `json:"limit_volume,omitempty" yaml:"limit_volume"`
	BlackStock     []string `json:"black_stock,omitempty" yaml:"black_stock"`
	BlackCategory  []string `json:"black_category,omitempty" yaml:"black_category"`
}

// Schedule Schedule
type Schedule struct {
	CleaneventCron     string `json:"cleanevent_cron,omitempty" yaml:"cleanevent_cron"`
	RestartSinopacCron string `json:"restart_sinopac_cron,omitempty" yaml:"restart_sinopac_cron"`
}

// parseConfig parseConfig
func parseConfig() {
	if globalConfig != nil {
		return
	}
	yamlFile, err := ioutil.ReadFile(global.Get().GetBasePath() + "/configs/config.yaml")
	if err != nil {
		log.Get().Panic(err)
	}
	err = yaml.Unmarshal(yamlFile, &globalConfig)
	if err != nil {
		log.Get().Panic(err)
	}

	// if in development, change some parameters
	if global.Get().GetIsDevelopment() {
		localHost := utils.GetHostIP()
		globalConfig.MQTT.Host = localHost

		globalConfig.Server.RunMode = "debug"
		globalConfig.Server.SinopacSRVHost = localHost

		globalConfig.Database.DBHost = localHost
		globalConfig.Database.Database = fmt.Sprintf("%s_debug", globalConfig.Database.Database)
	}

	globalConfig.basePath = global.Get().GetBasePath()
	globalConfig.MQTT.CAPath = globalConfig.basePath + globalConfig.MQTT.CAPath
	globalConfig.MQTT.KeyPath = globalConfig.basePath + globalConfig.MQTT.KeyPath
	globalConfig.MQTT.CertPath = globalConfig.basePath + globalConfig.MQTT.CertPath
}

// GetServerConfig GetServerConfig
func GetServerConfig() Server {
	if globalConfig != nil {
		return globalConfig.Server
	}
	once.Do(parseConfig)
	return globalConfig.Server
}

// GetDBConfig GetDBConfig
func GetDBConfig() Database {
	if globalConfig != nil {
		return globalConfig.Database
	}
	once.Do(parseConfig)
	return globalConfig.Database
}

// GetScheduleConfig GetScheduleConfig
func GetScheduleConfig() Schedule {
	if globalConfig != nil {
		return globalConfig.Schedule
	}
	once.Do(parseConfig)
	return globalConfig.Schedule
}

// GetTradeConfig GetTradeConfig
func GetTradeConfig() Trade {
	if globalConfig != nil {
		return globalConfig.Trade
	}
	once.Do(parseConfig)
	return globalConfig.Trade
}

// GetSwitchConfig GetSwitchConfig
func GetSwitchConfig() Switch {
	if globalConfig != nil {
		return globalConfig.Switch
	}
	once.Do(parseConfig)
	return globalConfig.Switch
}

// GetTargetCondConfig GetTargetCondConfig
func GetTargetCondConfig() TargetCond {
	if globalConfig != nil {
		return globalConfig.TargetCond
	}
	once.Do(parseConfig)
	return globalConfig.TargetCond
}

// GetMQConfig GetMQConfig
func GetMQConfig() MQTT {
	if globalConfig != nil {
		return globalConfig.MQTT
	}
	once.Do(parseConfig)
	return globalConfig.MQTT
}
