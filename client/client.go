package client

import (
	pb "FormSubmit/grpc"
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
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

func initializeConnection(host string, port uint) (*grpc.ClientConn, pb.FormSubmitClient) {

	conn, err := grpc.Dial(fmt.Sprintf("%v:%v", host, port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	return conn, pb.NewFormSubmitClient(conn)
}
func fillForm(args []string) *pb.FormData {
	form := pb.FormData{}
	fmt.Sscanf(args[0], "%s", &form.FirstName)
	fmt.Sscanf(args[1], "%s", &form.LastName)
	fmt.Sscanf(args[2], "%s", &form.Email)
	fmt.Sscanf(args[3], "%d", &form.Age)
	fmt.Sscanf(args[4], "%f", &form.Height)
	return &form
}
func viperSetup() (*defaultConfig, error) {
	var config *defaultConfig
	viper.SetConfigType("json")
	viper.SetConfigFile("./defaults.json")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	err = viper.Unmarshal(config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
func getHostAndPort(flags *pflag.FlagSet, config *defaultConfig) (string, uint) {
	host, _ := flags.GetString("host")
	if host == "" {
		host = config.DefaultHost
	}
	port, _ := flags.GetUint("port")
	if port == 0 {
		port = config.DefaultPort
	}
	return host, port
}
func RunClient(command *cobra.Command, args []string) {
	config, err := viperSetup()
	if err != nil {
		log.Fatal("Error reading config (defaults.json)! aborting...")
		return
	}

	// Contact the server and print out its response.
	host, port := getHostAndPort(command.Flags(), config)
	form := fillForm(args)
	conn, c := initializeConnection(host, port)
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SubmitForm(ctx, form)
	if err != nil {
		log.Printf("could not submit form: %v\n", err)
	} else {
		log.Printf("Status:%v, details:%v\n", r.GetSuccess(), r.GetDetails())
	}

}
