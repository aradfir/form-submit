package config

import (
	"github.com/spf13/viper"
)

type defaultConfig struct {
	DefaultHost string
	DefaultPort uint
}

func GetHostAndPort(host string, port uint, config *defaultConfig) (string, uint) {
	if host == "" {
		host = config.DefaultHost
	}
	if port == 0 {
		port = config.DefaultPort
	}
	return host, port
}

func ViperSetup() (defaultConfig, error) {
	var config defaultConfig
	viper.SetConfigType("json")
	viper.SetConfigFile("./configs/defaults.json")
	err := viper.ReadInConfig()
	if err != nil {
		return defaultConfig{}, err
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		return defaultConfig{}, err
	}
	return config, nil
}
