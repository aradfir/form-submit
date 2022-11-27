package server

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
)

func viperSetup() (defaultConfig, error) {
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

func setupMetricServer() {
	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(":2112", nil)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Error starting metrics server")
	}
}
