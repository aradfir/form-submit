package client

import (
	pb "FormSubmit/internal/grpc"
	"context"
	"fmt"
	"time"
)

func RunClient(data *pb.FormData, host string, port uint) {
	conn, c := initializeConnection(host, port)
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SubmitForm(ctx, data)
	if err != nil {
		fmt.Printf("could not submit form: %v\n", err)
	} else {
		fmt.Printf("Status:%v, details:%v\n", r.GetSuccess(), r.GetDetails())
	}

}
