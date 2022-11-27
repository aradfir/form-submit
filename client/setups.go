package client

import (
	pb "FormSubmit/grpc"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

func initializeConnection(host string, port uint) (*grpc.ClientConn, pb.FormSubmitClient) {

	conn, err := grpc.Dial(fmt.Sprintf("%v:%v", host, port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Fatal("could not connect to server")
	}

	return conn, pb.NewFormSubmitClient(conn)
}
