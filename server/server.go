package server

import (
	pb "FormSubmit/grpc"
	"context"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"math/rand"
	"net"
	"net/mail"
	"os"
	"reflect"
	"runtime"
	"time"
)

type defaultConfig struct {
	DefaultHost string
	DefaultPort uint
}

var config defaultConfig

type server struct {
	pb.UnimplementedFormSubmitServer
}

func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func EmailValidator(in *pb.FormData) bool {
	_, err := mail.ParseAddress(in.GetEmail())
	return err == nil
}
func AgeValidator(in *pb.FormData) bool {
	return in.GetAge() > 0
}
func HeightValidator(in *pb.FormData) bool {
	return in.GetHeight() > 0
}
func checkValidators(in *pb.FormData, validators []func(data *pb.FormData) bool) error {
	for _, validator := range validators {
		if !validator(in) {
			failedValidator := GetFunctionName(validator)
			errorText := fmt.Sprintf("validator %v failed", failedValidator)
			return errors.New(errorText)
		}
	}
	return nil
}
func safeClose(f *os.File) {
	err := f.Close()
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("File close error")
	}
}
func (s *server) SubmitForm(ctx context.Context, in *pb.FormData) (*pb.FormResult, error) {
	var validators = []func(data *pb.FormData) bool{EmailValidator, AgeValidator, HeightValidator}
	if err := checkValidators(in, validators); err != nil {
		return &pb.FormResult{
			Success: false,
			Details: err.Error(),
		}, err
	}

	f, err := os.OpenFile("users.form", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("DB file open failed")
		return &pb.FormResult{
			Success: false,
			Details: "Failed to open database",
		}, err
	}
	defer safeClose(f)

	_, err = f.WriteString(fmt.Sprintf("%v ### %v ### %v ### %v ### %v\n",
		in.GetFirstName(), in.GetLastName(), in.GetEmail(), in.GetAge(), in.GetHeight()))
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("DB file write failed")
		return &pb.FormResult{
			Success: false,
			Details: "Failed to write to server",
		}, err
	}

	return &pb.FormResult{
			Success: true,
			Details: "Hooray!"},
		nil
}
func loggerServerInterceptor(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {

	requestId := rand.Uint64()
	log.WithFields(log.Fields{"request ID": requestId}).Info("Request  started \n", requestId)
	start := time.Now()
	h, err := handler(ctx, req)
	log.WithFields(log.Fields{
		"request ID": requestId,
		"method":     info.FullMethod,
		"duration":   time.Since(start),
		"error":      err,
	}).Info("Request finished")
	return h, err

}

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
func getHostAndPort(host string, port uint, config *defaultConfig) (string, uint) {
	if host == "" {
		host = config.DefaultHost
	}
	if port == 0 {
		port = config.DefaultPort
	}
	return host, port
}
func RunServer(host string, port uint) {
	config, err := viperSetup()
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Error reading config")
	}
	host, port = getHostAndPort(host, port, &config)
	lis, err := net.Listen("tcp", fmt.Sprintf("%v:%v", host, port))
	if err != nil {
		log.WithFields(log.Fields{
			"host":  host,
			"port":  port,
			"error": err,
		}).Fatal("Failed to listen")
	}
	s := grpc.NewServer(grpc.UnaryInterceptor(loggerServerInterceptor))
	pb.RegisterFormSubmitServer(s, &server{})
	log.WithFields(log.Fields{"address": lis.Addr()}).Print("server listening")
	if err := s.Serve(lis); err != nil {
		log.WithFields(log.Fields{"error": err}).Fatal("failed to serve")
	}

}
