package cmd

import (
	"FormSubmit/internal/client"
	"FormSubmit/internal/config"
	pb "FormSubmit/internal/grpc"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Runs the client",
	Long: "Runs the client, you have to input" +
		" firstname, last name, email, age, height in order",
	Args: cobra.ExactArgs(5),
	Run:  runClient,
}

func init() {
	clientCmd.Flags().StringP("host", "H", "", "Host to send form to.")
	clientCmd.Flags().UintP("port", "P", 0, "Port to send form to.")
	RootCmd.AddCommand(clientCmd)
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
func runClient(command *cobra.Command, args []string) {
	cfg, err := config.ViperSetup()
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Error reading config")
		return
	}
	hostFlag, err := command.Flags().GetString("host")
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("bad host flag")
		return
	}
	portFlag, err := command.Flags().GetUint("port")
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("bad port flag")
		return
	}
	host, port := config.GetHostAndPort(hostFlag, portFlag, &cfg)
	client.RunClient(fillForm(args), host, port)
}
