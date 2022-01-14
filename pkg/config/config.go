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
	basePath   string
	Server     Server     `json:"server,omitempty" yaml:"server"`
	Database   Database   `json:"database,omitempty" yaml:"database"`
	MQTT       MQTT       `json:"mqtt,omitempty" yaml:"mqtt"`
	Switch     Switch     `json:"switch,omitempty" yaml:"switch"`
	Trade      Trade      `json:"trade,omitempty" yaml:"trade"`
	Quota      Quota      `json:"quota,omitempty" yaml:"quota"`
	TargetCond TargetCond `json:"target_cond,omitempty" yaml:"target_cond"`
	Analyze    Analyze    `json:"analyze,omitempty" yaml:"analyze"`
	Schedule   Schedule   `json:"schedule,omitempty" yaml:"schedule"`
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
	Host     string `json:"host,omitempty" yaml:"host"`
	Port     string `json:"port,omitempty" yaml:"port"`
	User     string `json:"user,omitempty" yaml:"user"`
	Passwd   string `json:"passwd,omitempty" yaml:"passwd"`
	Database string `json:"database,omitempty" yaml:"database"`
	TimeZone string `json:"time_zone,omitempty" yaml:"time_zone"`
}

// MQTT MQTT
type MQTT struct {
	Host     string `json:"host,omitempty" yaml:"host"`
	Port     string `json:"port,omitempty" yaml:"port"`
	User     string `json:"user,omitempty" yaml:"user"`
	Passwd   string `json:"passwd,omitempty" yaml:"passwd"`
	ClientID string `json:"client_id,omitempty" yaml:"client_id"`
	CAPath   string `json:"ca_path,omitempty" yaml:"ca_path"`
	CertPath string `json:"cert_path,omitempty" yaml:"cert_path"`
	KeyPath  string `json:"key_path,omitempty" yaml:"key_path"`
}

// Switch Switch
type Switch struct {
	Simulation bool `json:"simulation,omitempty" yaml:"simulation"`

	Buy       bool `json:"buy,omitempty" yaml:"buy"`
	Sell      bool `json:"sell,omitempty" yaml:"sell"`
	SellFirst bool `json:"sell_first,omitempty" yaml:"sell_first"`
	BuyLater  bool `json:"buy_later,omitempty" yaml:"buy_later"`

	MeanTimeForward int64 `json:"mean_time_forward,omitempty" yaml:"mean_time_forward"`
	MeanTimeReverse int64 `json:"mean_time_reverse,omitempty" yaml:"mean_time_reverse"`
	ForwardMax      int64 `json:"forward_max,omitempty" yaml:"forward_max"`
	ReverseMax      int64 `json:"reverse_max,omitempty" yaml:"reverse_max"`
}

// Trade Trade
type Trade struct {
	HistoryClosePeriod int64 `json:"history_close_period,omitempty" yaml:"history_close_period"`
	HistoryTickPeriod  int64 `json:"history_tick_period,omitempty" yaml:"history_tick_period"`
	HistoryKbarPeriod  int64 `json:"history_kbar_period,omitempty" yaml:"history_kbar_period"`

	HoldTimeFromOpen float64 `json:"hold_time_from_open,omitempty" yaml:"hold_time_from_open"`
	TotalOpenTime    float64 `json:"total_open_time,omitempty" yaml:"total_open_time"`

	TradeInWaitTime  int64 `json:"trade_in_wait_time,omitempty" yaml:"trade_in_wait_time"`
	TradeOutWaitTime int64 `json:"trade_out_wait_time,omitempty" yaml:"trade_out_wait_time"`

	TradeInEndTime  float64 `json:"trade_in_end_time,omitempty" yaml:"trade_in_end_time"`
	TradeOutEndTime float64 `json:"trade_out_end_time,omitempty" yaml:"trade_out_end_time"`
}

// Quota Quota
type Quota struct {
	TradeQuota    int64   `json:"trade_quota,omitempty" yaml:"trade_quota"`
	TradeTaxRatio float64 `json:"trade_tax_ratio,omitempty" yaml:"trade_tax_ratio"`
	TradeFeeRatio float64 `json:"trade_fee_ratio,omitempty" yaml:"trade_fee_ratio"`
	FeeDiscount   float64 `json:"fee_discount,omitempty" yaml:"fee_discount"`
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

// Analyze Analyze
type Analyze struct {
	CloseChangeRatioLow  float64 `json:"close_change_ratio_low,omitempty" yaml:"close_change_ratio_low"`
	CloseChangeRatioHigh float64 `json:"close_change_ratio_high,omitempty" yaml:"close_change_ratio_high"`
	OpenCloseChangeRatio float64 `json:"open_close_change_ratio,omitempty" yaml:"open_close_change_ratio"`
	OutInRatio           float64 `json:"out_in_ratio,omitempty" yaml:"out_in_ratio"`
	InOutRatio           float64 `json:"in_out_ratio,omitempty" yaml:"in_out_ratio"`

	VolumePRLow  float64 `json:"volume_pr_low,omitempty" yaml:"volume_pr_low"`
	VolumePRHigh float64 `json:"volume_pr_high,omitempty" yaml:"volume_pr_high"`

	TickAnalyzeMinPeriod float64 `json:"tick_analyze_min_period,omitempty" yaml:"tick_analyze_min_period"`
	TickAnalyzeMaxPeriod float64 `json:"tick_analyze_max_period,omitempty" yaml:"tick_analyze_max_period"`

	RSIMinCount int `json:"rsi_min_count,omitempty" yaml:"rsi_min_count"`

	RSIHigh float64 `json:"rsi_high,omitempty" yaml:"rsi_high"`
	RSILow  float64 `json:"rsi_low,omitempty" yaml:"rsi_low"`

	MaxLoss float64 `json:"max_loss,omitempty" yaml:"max_loss"`
}

// Schedule Schedule
type Schedule struct {
	CleanEvent     string `json:"clean_event,omitempty" yaml:"clean_event"`
	RestartSinopac string `json:"restart_sinopac,omitempty" yaml:"restart_sinopac"`
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

		globalConfig.Switch.Simulation = true
	}

	globalConfig.basePath = basePath
	globalConfig.MQTT.CAPath = filepath.Join(globalConfig.basePath, globalConfig.MQTT.CAPath)
	globalConfig.MQTT.KeyPath = filepath.Join(globalConfig.basePath, globalConfig.MQTT.KeyPath)
	globalConfig.MQTT.CertPath = filepath.Join(globalConfig.basePath, globalConfig.MQTT.CertPath)

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
func GetSwitchConfig() Switch {
	if globalConfig != nil {
		return globalConfig.Switch
	}
	once.Do(parseConfig)
	return globalConfig.Switch
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
