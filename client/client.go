package client

import (
	pb "FormSubmit/grpc"
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

func initalizeConnection(host string, port uint) (*grpc.ClientConn, pb.FormSubmitClient) {

	conn, err := grpc.Dial(fmt.Sprintf("%v:%v", host, port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	return conn, pb.NewFormSubmitClient(conn)
}
func fillForm(args []string) *pb.FormData {
	form := pb.FormData{}
	//fmt.Println("Enter your First Name:")
	fmt.Sscanf(args[0], "%s", &form.FirstName)
	//fmt.Println("Enter your Last Name:")
	fmt.Sscanf(args[1], "%s", &form.LastName)
	//fmt.Println("Enter your Email:")
	fmt.Sscanf(args[2], "%s", &form.Email)
	//fmt.Println("Enter your Age:")
	fmt.Sscanf(args[3], "%d", &form.Age)
	//fmt.Println("Enter your Height (m):")
	fmt.Sscanf(args[4], "%f", &form.Height)
	return &form
}

func RunClient(command *cobra.Command, args []string) {

	// Contact the server and print out its response.
	host, _ := command.Flags().GetString("host")
	port, _ := command.Flags().GetUint("port")
	form := fillForm(args)
	conn, c := initalizeConnection(host, port)
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	r, err := c.SubmitForm(ctx, form)
	if err != nil {
		log.Printf("could not submit form: %v\n", err)
	} else {
		log.Printf("Status:%v, details:%v\n", r.GetSuccess(), r.GetDetails())
	}
	defer cancel()
}
