package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/levelord1311/backendForSharedProject/lot_service/pkg/logging"
	"sync"
)

type Config struct {
	IsDebug *bool `yaml:"is_debug"`
	Listen  struct {
		Type   string `yaml:"type" env-default:"port"`
		BindIP string `yaml:"bind_ip" env-default:"localhost"`
		Port   string `yaml:"port" env-default:"8082"`
	}

	MysqlDB struct {
		Host     string `yaml:"host" env-default:"0.0.0.0"`
		Port     string `yaml:"port" env-default:"3306"`
		Username string `yaml:"username" env-default:"testDB"`
		Password string `yaml:"password" env-default:"testPassword"`
		DBName   string `yaml:"db_name" env-default:"test_db"`
	}
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		logger := logging.GetLogger()
		logger.Info("read application config")
		instance = &Config{}
		if err := cleanenv.ReadConfig("config.yml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Info(help)
			logger.Fatal(err)
		}
	})
	return instance
}
