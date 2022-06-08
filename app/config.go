package app

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	Database DatabaseConfig `mapstructure:"database"`
	Server   ServerConfig   `mapstructure:"server"`
}

type DatabaseConfig struct {
	Driver   string `mapstructure:"driver"`
	Domain   string `mapstructure:"domain"`
	Port     string `mapstructure:"port"`
	Name     string `mapstructure:"name"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type ServerConfig struct {
	Address string `mapstructure:"address"`
	Port    string `mapstructure:"port"`
}

func NewDatabaseConfig() DatabaseConfig {
	return unmarshalConfig().Database
}

func NewServerConfig() ServerConfig {
	return unmarshalConfig().Server
}

func unmarshalConfig() Config {
	var config Config
	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("json")
	v.AddConfigPath(".")

	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}
	log.Debug("Database config loaded")

	err = v.Unmarshal(&config)
	if err != nil {
		panic(fmt.Sprintf("Unable to decode into struct, %v", err))
	}
	log.Debug("Database config unmarshalled")

	return config
}
