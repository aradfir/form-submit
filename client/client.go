package main

import (
	pb "FormSubmit/grpc"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

const (
	addr = "localhost:8080"
)

func initalizeConnection() (*grpc.ClientConn, pb.FormSubmitClient) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	return conn, pb.NewFormSubmitClient(conn)
}
func fillForm() *pb.FormData {
	form := pb.FormData{}
	fmt.Println("Enter your First Name:")
	fmt.Scanf("%s", &form.FirstName)
	fmt.Println("Enter your Last Name:")
	fmt.Scanf("%s", &form.LastName)
	fmt.Println("Enter your Email:")
	fmt.Scanf("%s", &form.Email)
	fmt.Println("Enter your Age:")
	fmt.Scanf("%d", &form.Age)
	fmt.Println("Enter your Height (m):")
	fmt.Scanf("%f", &form.Height)
	return &form
}
func main() {

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	conn, c := initalizeConnection()
	defer conn.Close()
	anotherUser := "y"
	for anotherUser == "y" {
		form := fillForm()
		ctx, cancel = context.WithTimeout(context.Background(), time.Second)
		r, err := c.SubmitForm(ctx, form)
		if err != nil {
			log.Printf("could not submit form: %v\n", err)
		} else {
			log.Printf("Status:%v, details:%v\n", r.GetSuccess(), r.GetDetails())
		}
		fmt.Println("Would yo like to submit another user? (y/[n])")
		fmt.Scanf("%s", &anotherUser)

	}
	defer cancel()
}
