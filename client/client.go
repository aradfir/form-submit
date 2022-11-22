package client

import (
	pb "FormSubmit/grpc"
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

type defaultConfig struct {
	DefaultHost string
	DefaultPort uint
}

var config defaultConfig

func initializeConnection(host string, port uint) (*grpc.ClientConn, pb.FormSubmitClient) {

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
	// Contact the server and print out its response.
	host, _ := command.Flags().GetString("host")
	if host == "" {
		host = config.DefaultHost
	}
	port, _ := command.Flags().GetUint("port")
	if port == 0 {
		port = config.DefaultPort
	}
	form := fillForm(args)
	conn, c := initializeConnection(host, port)
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
