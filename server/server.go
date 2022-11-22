package server

import (
	pb "FormSubmit/grpc"
	"context"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"log"
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
func (s *server) SubmitForm(ctx context.Context, in *pb.FormData) (*pb.FormResult, error) {
	var validators = [...]func(data *pb.FormData) bool{EmailValidator, AgeValidator, HeightValidator}
	for _, validator := range validators {
		var correct bool = validator(in)
		if !correct {
			failedValidator := GetFunctionName(validator)
			errorText := fmt.Sprintf("validator %v failed", failedValidator)
			return &pb.FormResult{
				Success: false,
				Details: errorText,
			}, errors.New(errorText)
		}
	}

	f, err := os.OpenFile("users.form", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return &pb.FormResult{
			Success: false,
			Details: "Failed to open database",
		}, err
	}

	defer func() {
		err := f.Close()
		if err != nil {
			log.Fatalf("File close error:%v", err)
		}
	}()
	_, err = f.WriteString(fmt.Sprintf("%v ### %v ### %v ### %v ### %v\n", in.GetFirstName(), in.GetLastName(), in.GetEmail(), in.GetAge(), in.GetHeight()))
	if err != nil {
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
	log.Printf("Request %v started \n", requestId)
	start := time.Now()
	h, err := handler(ctx, req)
	log.Printf("Request %v finished - method:%s\tduration:%s\tError:%v\n", requestId, info.FullMethod, time.Since(start), err)
	return h, err

}
func RunServer(address string, port uint) {
	viper.SetConfigType("json")
	viper.SetConfigFile("./defaults.json")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("Error reading default config! aborting...")
		return
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		fmt.Println("Error parsing default config! aborting...")
		return
	}
	if port == 0 {
		port = config.DefaultPort
	}
	if address == "" {
		address = config.DefaultHost
	}
	lis, err := net.Listen("tcp", fmt.Sprintf("%v:%v", address, port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer(grpc.UnaryInterceptor(loggerServerInterceptor))
	pb.RegisterFormSubmitServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
