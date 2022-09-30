package config

import (
	"github.com/spf13/viper"
)

// PsqlConfig stores all configuration for postgres DB.
// Variables are read from .env file
type PsqlConfig struct {
	Host     string `mapstucture:"HOST"`
	Port     int    `mapstucture:"PORT"`
	User     string `mapstucture:"USER"`
	Password string `mapstucture:"PASSWORD"`
	DBName   string `mapstucture:"DBNAME"`
}

type JWTConfig struct {
	Key string `mapstructure:"KEY"`
}

type Config struct {
	Path string `mapstructure:"CONFIG_PATH"`
}

func LoadMainConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}

// LoadPSQLConfig reads configuration from file or environment variables.
func LoadPSQLConfig(path string) (config PsqlConfig, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("psql")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}

func LoadJWTConfig(path string) (config JWTConfig, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("jwt")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
