// Package config package config
package config

import (
	"io/ioutil"
	"sync"

	"gitlab.tocraw.com/root/toc_trader/pkg/log"
	"gopkg.in/yaml.v2"
)

var (
	globalConfig *Config
	initLock     sync.Mutex
)

// Config Config
type Config struct {
	Server   Server   `yaml:"server"`
	Database Database `yaml:"database"`
}

// Server Server
type Server struct {
	RunMode        string `yaml:"run_mode"`
	HTTPPort       string `yaml:"http_port"`
	SinopacSRVHost string `yaml:"sinopac_srv_host"`
	SinopacSRVPort string `yaml:"sinopac_srv_port"`
	Reset          bool   `yaml:"reset"`
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

// parseConfig parseConfig
func parseConfig() (err error) {
	defer initLock.Unlock()
	initLock.Lock()
	if globalConfig != nil {
		return
	}

	var yamlFile []byte
	yamlFile, err = ioutil.ReadFile("./configs/config.yaml")
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(yamlFile, &globalConfig)
	if err != nil {
		return err
	}
	return err
}

// Get Get
func Get() (config Config, err error) {
	if globalConfig == nil {
		err = parseConfig()
		if err != nil {
			log.Get().Panic(err)
		}
	}
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
