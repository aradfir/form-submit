package server

import (
	"FormSubmit/internal/grpc"
	"FormSubmit/internal/interceptors"
	"context"
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"os"
)

type Validator func(data *form_data.FormData) bool
type server struct {
	form_data.UnimplementedFormSubmitServer
	formSubmitValidators []Validator
}

func safeClose(f *os.File) {
	err := f.Close()
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("File close error")
	}
}
func (s *server) SubmitForm(ctx context.Context, in *form_data.FormData) (*form_data.FormResult, error) {

	if err := checkValidators(in, s.formSubmitValidators); err != nil {
		return &form_data.FormResult{
			Success: false,
			Details: err.Error(),
		}, err
	}

	f, err := os.OpenFile("users.form", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("DB file open failed")
		return &form_data.FormResult{
			Success: false,
			Details: "Failed to open database",
		}, err
	}
	defer safeClose(f)

	_, err = f.WriteString(fmt.Sprintf("%v ### %v ### %v ### %v ### %v\n",
		in.GetFirstName(), in.GetLastName(), in.GetEmail(), in.GetAge(), in.GetHeight()))
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("DB file write failed")
		return &form_data.FormResult{
			Success: false,
			Details: "Failed to write to server",
		}, err
	}

	return &form_data.FormResult{
			Success: true,
			Details: "Hooray!"},
		nil
}

func RunServer(host string, port uint) {
	lis, err := net.Listen("tcp", fmt.Sprintf("%v:%v", host, port))
	if err != nil {
		log.WithFields(log.Fields{
			"host": host,
			"port": port,

			"error": err,
		}).Fatal("Failed to listen")
	}
	s := grpc.NewServer(grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(interceptors.ServerLogger, interceptors.ServerInstrumentalization)))
	form_data.RegisterFormSubmitServer(s, &server{
		formSubmitValidators: []Validator{EmailValidator, AgeValidator, HeightValidator},
	})
	log.WithFields(log.Fields{"address": lis.Addr()}).Print("server listening")
	if err := s.Serve(lis); err != nil {
		log.WithFields(log.Fields{"error": err}).Fatal("failed to serve")
	}

}
