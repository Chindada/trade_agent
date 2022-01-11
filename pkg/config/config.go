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
	Analyze    Analyze    `json:"analyze,omitempty" yaml:"analyze"`
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

	TradeInWaitTime  int64 `json:"trade_in_wait_time,omitempty" yaml:"trade_in_wait_time"`
	TradeOutWaitTime int64 `json:"trade_out_wait_time,omitempty" yaml:"trade_out_wait_time"`

	WaitInOpen      int64 `json:"wait_in_open,omitempty" yaml:"wait_in_open"`
	TradeInEndTime  int64 `json:"trade_in_end_time,omitempty" yaml:"trade_in_end_time"`
	TradeOutEndTime int64 `json:"trade_out_end_time,omitempty" yaml:"trade_out_end_time"`
}

// Switch Switch
type Switch struct {
	EnableBuy       bool `json:"enable_buy,omitempty" yaml:"enable_buy"`
	EnableSell      bool `json:"enable_sell,omitempty" yaml:"enable_sell"`
	EnableSellFirst bool `json:"enable_sell_first,omitempty" yaml:"enable_sell_first"`
	EnableBuyLater  bool `json:"enable_buy_later,omitempty" yaml:"enable_buy_later"`

	MeanTimeForward int64 `json:"mean_time_forward,omitempty" yaml:"mean_time_forward"`
	MeanTimeReverse int64 `json:"mean_time_reverse,omitempty" yaml:"mean_time_reverse"`
	ForwardMax      int64 `json:"forward_max,omitempty" yaml:"forward_max"`
	ReverseMax      int64 `json:"reverse_max,omitempty" yaml:"reverse_max"`
}

// TargetCond TargetCond
type TargetCond struct {
	LimitPriceLow        float64  `json:"limit_price_low,omitempty" yaml:"limit_price_low"`
	LimitPriceHigh       float64  `json:"limit_price_high,omitempty" yaml:"limit_price_high"`
	LimitVolume          int64    `json:"limit_volume,omitempty" yaml:"limit_volume"`
	BlackStock           []string `json:"black_stock,omitempty" yaml:"black_stock"`
	BlackCategory        []string `json:"black_category,omitempty" yaml:"black_category"`
	RealTimeTargetsCount int64    `json:"real_time_targets_count,omitempty" yaml:"real_time_targets_count"`
}

// Schedule Schedule
type Schedule struct {
	CleanEvent     string `json:"clean_event,omitempty" yaml:"clean_event"`
	RestartSinopac string `json:"restart_sinopac,omitempty" yaml:"restart_sinopac"`
}

// Analyze Analyze
type Analyze struct {
	CloseChangeRatioLow  float64 `json:"close_change_ratio_low,omitempty" yaml:"close_change_ratio_low"`
	CloseChangeRatioHigh float64 `json:"close_change_ratio_high,omitempty" yaml:"close_change_ratio_high"`
	OpenCloseChangeRatio float64 `json:"open_close_change_ratio,omitempty" yaml:"open_close_change_ratio"`
	OutInRatio           float64 `json:"out_in_ratio,omitempty" yaml:"out_in_ratio"`
	InOutRatio           float64 `json:"in_out_ratio,omitempty" yaml:"in_out_ratio"`
	VolumePR             float64 `json:"volume_pr,omitempty" yaml:"volume_pr"`

	TickAnalyzeMinPeriod int64 `json:"tick_analyze_period,omitempty" yaml:"tick_analyze_period"`
	TickAnalyzeMaxPeriod int64 `json:"tick_analyze_max_period,omitempty" yaml:"tick_analyze_max_period"`

	RSIMinCount int `json:"rsi_min_count,omitempty" yaml:"rsi_min_count"`

	RSIHigh float64 `json:"rsi_high,omitempty" yaml:"rsi_high"`
	RSILow  float64 `json:"rsi_low,omitempty" yaml:"rsi_low"`
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

		globalConfig.Switch.EnableBuy = false
		globalConfig.Switch.EnableSellFirst = false
	}

	globalConfig.basePath = global.Get().GetBasePath()
	globalConfig.MQTT.CAPath = globalConfig.basePath + globalConfig.MQTT.CAPath
	globalConfig.MQTT.KeyPath = globalConfig.basePath + globalConfig.MQTT.KeyPath
	globalConfig.MQTT.CertPath = globalConfig.basePath + globalConfig.MQTT.CertPath

	checkConfigIsEmpty(*globalConfig)
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

// GetAnalyzeConfig GetAnalyzeConfig
func GetAnalyzeConfig() Analyze {
	if globalConfig != nil {
		return globalConfig.Analyze
	}
	once.Do(parseConfig)
	return globalConfig.Analyze
}
