package server

import (
	pb "FormSubmit/grpc"
	"context"
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"os"
)

type defaultConfig struct {
	DefaultHost string
	DefaultPort uint
}

type server struct {
	pb.UnimplementedFormSubmitServer
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
	go setupMetricServer()
	host, port = getHostAndPort(host, port, &config)
	lis, err := net.Listen("tcp", fmt.Sprintf("%v:%v", host, port))
	if err != nil {
		log.WithFields(log.Fields{
			"host":  host,
			"port":  port,
			"error": err,
		}).Fatal("Failed to listen")
	}
	s := grpc.NewServer(grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(loggerServerInterceptor, instrumentalizationInterceptor)))

	pb.RegisterFormSubmitServer(s, &server{})
	log.WithFields(log.Fields{"address": lis.Addr()}).Print("server listening")
	if err := s.Serve(lis); err != nil {
		log.WithFields(log.Fields{"error": err}).Fatal("failed to serve")
	}

}
