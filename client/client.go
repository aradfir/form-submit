package client

import (
	pb "FormSubmit/grpc"
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"time"
)

type defaultConfig struct {
	DefaultHost string
	DefaultPort uint
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
		log.WithFields(log.Fields{"error": err}).Error("Error reading config")
		return
	}
	host, port := getHostAndPort(command.Flags(), &config)
	form := fillForm(args)
	conn, c := initializeConnection(host, port)
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SubmitForm(ctx, form)
	if err != nil {
		fmt.Printf("could not submit form: %v\n", err)
	} else {
		fmt.Printf("Status:%v, details:%v\n", r.GetSuccess(), r.GetDetails())
	}

}
