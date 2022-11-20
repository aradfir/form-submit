package server

import (
	pb "FormSubmit/grpc"
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/mail"
	"os"
	"reflect"
	"runtime"
	"strconv"
)

var port = 8080

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
func RunServer() {
	lis, err := net.Listen("tcp", "localhost:"+strconv.Itoa(port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterFormSubmitServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
