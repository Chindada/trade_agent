// Package config package config
package config

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
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
	basePath    string
	Server      Server      `json:"server" yaml:"server"`
	Database    Database    `json:"database" yaml:"database"`
	MQTT        MQTT        `json:"mqtt" yaml:"mqtt"`
	TradeSwitch TradeSwitch `json:"trade_switch" yaml:"trade_switch"`
	Trade       Trade       `json:"trade" yaml:"trade"`
	Quota       Quota       `json:"quota" yaml:"quota"`
	TargetCond  TargetCond  `json:"target_cond" yaml:"target_cond"`
	Analyze     Analyze     `json:"analyze" yaml:"analyze"`
	Schedule    Schedule    `json:"schedule" yaml:"schedule"`
}

// Server Server
type Server struct {
	RunMode        string `json:"run_mode" yaml:"run_mode"`
	HTTPPort       string `json:"http_port" yaml:"http_port"`
	SinopacSRVHost string `json:"sinopac_srv_host" yaml:"sinopac_srv_host"`
	SinopacSRVPort string `json:"sinopac_srv_port" yaml:"sinopac_srv_port"`
	SinopacMAXConn int    `json:"sinopac_max_conn" yaml:"sinopac_max_conn"`
}

// Database Database
type Database struct {
	Host     string `json:"host" yaml:"host"`
	Port     string `json:"port" yaml:"port"`
	User     string `json:"user" yaml:"user"`
	Passwd   string `json:"passwd" yaml:"passwd"`
	Database string `json:"database" yaml:"database"`
	TimeZone string `json:"time_zone" yaml:"time_zone"`
}

// MQTT MQTT
type MQTT struct {
	Host     string `json:"host" yaml:"host"`
	Port     string `json:"port" yaml:"port"`
	User     string `json:"user" yaml:"user"`
	Passwd   string `json:"passwd" yaml:"passwd"`
	ClientID string `json:"client_id" yaml:"client_id"`
	CAPath   string `json:"ca_path" yaml:"ca_path"`
	CertPath string `json:"cert_path" yaml:"cert_path"`
	KeyPath  string `json:"key_path" yaml:"key_path"`
}

// TradeSwitch TradeSwitch
type TradeSwitch struct {
	Simulation bool `json:"simulation" yaml:"simulation"`

	Buy       bool `json:"buy" yaml:"buy"`
	Sell      bool `json:"sell" yaml:"sell"`
	SellFirst bool `json:"sell_first" yaml:"sell_first"`
	BuyLater  bool `json:"buy_later" yaml:"buy_later"`

	MeanTimeForward int64 `json:"mean_time_forward" yaml:"mean_time_forward"`
	MeanTimeReverse int64 `json:"mean_time_reverse" yaml:"mean_time_reverse"`
	ForwardMax      int64 `json:"forward_max" yaml:"forward_max"`
	ReverseMax      int64 `json:"reverse_max" yaml:"reverse_max"`
}

// Trade Trade
type Trade struct {
	HistoryClosePeriod int64 `json:"history_close_period" yaml:"history_close_period"`
	HistoryTickPeriod  int64 `json:"history_tick_period" yaml:"history_tick_period"`
	HistoryKbarPeriod  int64 `json:"history_kbar_period" yaml:"history_kbar_period"`

	HoldTimeFromOpen float64 `json:"hold_time_from_open" yaml:"hold_time_from_open"`
	TotalOpenTime    float64 `json:"total_open_time" yaml:"total_open_time"`

	TradeInWaitTime  int64 `json:"trade_in_wait_time" yaml:"trade_in_wait_time"`
	TradeOutWaitTime int64 `json:"trade_out_wait_time" yaml:"trade_out_wait_time"`

	TradeInEndTime  float64 `json:"trade_in_end_time" yaml:"trade_in_end_time"`
	TradeOutEndTime float64 `json:"trade_out_end_time" yaml:"trade_out_end_time"`
}

// Quota Quota
type Quota struct {
	TradeQuota    int64   `json:"trade_quota" yaml:"trade_quota"`
	TradeTaxRatio float64 `json:"trade_tax_ratio" yaml:"trade_tax_ratio"`
	TradeFeeRatio float64 `json:"trade_fee_ratio" yaml:"trade_fee_ratio"`
	FeeDiscount   float64 `json:"fee_discount" yaml:"fee_discount"`
}

// TargetCond TargetCond
type TargetCond struct {
	LimitPriceLow        float64  `json:"limit_price_low" yaml:"limit_price_low"`
	LimitPriceHigh       float64  `json:"limit_price_high" yaml:"limit_price_high"`
	LimitVolume          int64    `json:"limit_volume" yaml:"limit_volume"`
	BlackStock           []string `json:"black_stock" yaml:"black_stock"`
	BlackCategory        []string `json:"black_category" yaml:"black_category"`
	RealTimeTargetsCount int64    `json:"real_time_targets_count" yaml:"real_time_targets_count"`
}

// Analyze Analyze
type Analyze struct {
	CloseChangeRatioLow  float64 `json:"close_change_ratio_low" yaml:"close_change_ratio_low"`
	CloseChangeRatioHigh float64 `json:"close_change_ratio_high" yaml:"close_change_ratio_high"`

	OpenCloseChangeRatioLow  float64 `json:"open_close_change_ratio_low" yaml:"open_close_change_ratio_low"`
	OpenCloseChangeRatioHigh float64 `json:"open_close_change_ratio_high" yaml:"open_close_change_ratio_high"`

	OutInRatio float64 `json:"out_in_ratio" yaml:"out_in_ratio"`
	InOutRatio float64 `json:"in_out_ratio" yaml:"in_out_ratio"`

	VolumePRLow  float64 `json:"volume_pr_low" yaml:"volume_pr_low"`
	VolumePRHigh float64 `json:"volume_pr_high" yaml:"volume_pr_high"`

	TickAnalyzeMinPeriod float64 `json:"tick_analyze_min_period" yaml:"tick_analyze_min_period"`
	TickAnalyzeMaxPeriod float64 `json:"tick_analyze_max_period" yaml:"tick_analyze_max_period"`

	RSIMinCount int     `json:"rsi_min_count" yaml:"rsi_min_count"`
	RSIHigh     float64 `json:"rsi_high" yaml:"rsi_high"`
	RSILow      float64 `json:"rsi_low" yaml:"rsi_low"`

	MaxLoss  float64 `json:"max_loss" yaml:"max_loss"`
	MAPeriod int64   `json:"ma_period" yaml:"ma_period"`
}

// Schedule Schedule
type Schedule struct {
	CleanEvent     string `json:"clean_event" yaml:"clean_event"`
	RestartSinopac string `json:"restart_sinopac" yaml:"restart_sinopac"`
}

// parseConfig parseConfig
func parseConfig() {
	if globalConfig != nil {
		return
	}

	basePath := global.Get().GetBasePath()
	if basePath == "" {
		basePath = global.GetRuntimePath()
	}

	yamlFile, err := ioutil.ReadFile(filepath.Clean(filepath.Join(basePath, "/configs/config.yaml")))
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

		globalConfig.Database.Host = localHost
		globalConfig.Database.Database = fmt.Sprintf("%s_debug", globalConfig.Database.Database)

		globalConfig.TradeSwitch.Simulation = true
		globalConfig.TradeSwitch.Buy = false
		globalConfig.TradeSwitch.SellFirst = false
	}

	globalConfig.basePath = basePath
	globalConfig.MQTT.CAPath = filepath.Join(globalConfig.basePath, globalConfig.MQTT.CAPath)
	globalConfig.MQTT.KeyPath = filepath.Join(globalConfig.basePath, globalConfig.MQTT.KeyPath)
	globalConfig.MQTT.CertPath = filepath.Join(globalConfig.basePath, globalConfig.MQTT.CertPath)

	checkConfigIsEmpty(*globalConfig)
}

// GetAllConfig GetAllConfig
func GetAllConfig() Config {
	if globalConfig != nil {
		return *globalConfig
	}
	once.Do(parseConfig)
	return *globalConfig
}

// GetServerConfig GetServerConfig
func GetServerConfig() Server {
	if globalConfig != nil {
		return globalConfig.Server
	}
	once.Do(parseConfig)
	return globalConfig.Server
}

// GetDatabaseConfig GetDatabaseConfig
func GetDatabaseConfig() Database {
	if globalConfig != nil {
		return globalConfig.Database
	}
	once.Do(parseConfig)
	return globalConfig.Database
}

// GetMQTTConfig GetMQTTConfig
func GetMQTTConfig() MQTT {
	if globalConfig != nil {
		return globalConfig.MQTT
	}
	once.Do(parseConfig)
	return globalConfig.MQTT
}

// GetSwitchConfig GetSwitchConfig
func GetSwitchConfig() TradeSwitch {
	if globalConfig != nil {
		return globalConfig.TradeSwitch
	}
	once.Do(parseConfig)
	return globalConfig.TradeSwitch
}

// GetTradeConfig GetTradeConfig
func GetTradeConfig() Trade {
	if globalConfig != nil {
		return globalConfig.Trade
	}
	once.Do(parseConfig)
	return globalConfig.Trade
}

// GetQuotaConfig GetQuotaConfig
func GetQuotaConfig() Quota {
	if globalConfig != nil {
		return globalConfig.Quota
	}
	once.Do(parseConfig)
	return globalConfig.Quota
}

// GetTargetCondConfig GetTargetCondConfig
func GetTargetCondConfig() TargetCond {
	if globalConfig != nil {
		return globalConfig.TargetCond
	}
	once.Do(parseConfig)
	return globalConfig.TargetCond
}

// GetAnalyzeConfig GetAnalyzeConfig
func GetAnalyzeConfig() Analyze {
	if globalConfig != nil {
		return globalConfig.Analyze
	}
	once.Do(parseConfig)
	return globalConfig.Analyze
}

// GetScheduleConfig GetScheduleConfig
func GetScheduleConfig() Schedule {
	if globalConfig != nil {
		return globalConfig.Schedule
	}
	once.Do(parseConfig)
	return globalConfig.Schedule
}
