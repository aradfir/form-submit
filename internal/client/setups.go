package client

import (
	pb "FormSubmit/internal/grpc"
	"fmt"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func initializeConnection(host string, port uint) (*grpc.ClientConn, pb.FormSubmitClient) {

	conn, err := grpc.Dial(fmt.Sprintf("%v:%v", host, port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Fatal("could not connect to server")
	}

	return conn, pb.NewFormSubmitClient(conn)
}
