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
)

type Validator func(data *form_data.FormData) bool
type server struct {
	form_data.UnimplementedFormSubmitServer
	formSubmitValidators []Validator
}

func (s *server) SubmitForm(ctx context.Context, in *form_data.FormData) (*form_data.FormResult, error) {
	if err := checkValidators(in, s.formSubmitValidators); err != nil {
		return &form_data.FormResult{
			Success: false,
			Details: err.Error(),
		}, err
	}

	if err := WriteToFile(in); err != nil {
		return &form_data.FormResult{
				Success: false,
				Details: "Error writing to database",
			},
			err
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
