package main

import (
	"FormSubmit/client"
	"FormSubmit/server"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "FormSubmit",
	Short: "Form submit is a sample project to work with go, cobra, viper and grpc",
	Run: func(cmd *cobra.Command, args []string) {

	},
}
var commit string
var buildTime string

func main() {

	clientCommand := &cobra.Command{
		Use:   "client",
		Short: "Runs the client",
		Long: "Runs the client, you have to input" +
			" firstname, last name, email, age, height in order",
		Args: cobra.ExactArgs(5),
		Run:  client.RunClient,
	}
	var (
		address string
		port    uint
	)
	serverCommand := &cobra.Command{
		Use:   "serve",
		Short: "Runs the server",
		Args:  cobra.NoArgs,

		Run: func(cmd *cobra.Command, args []string) {
			server.RunServer(address, port)
		},
	}
	versionCommand := &cobra.Command{
		Use:   "version",
		Short: "Get the version of this program",
		Long:  "Info about the program aka build sha and build time",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Commit SHA:", commit)
			fmt.Println("Build time:", buildTime)
		},
	}
	serverCommand.Flags().StringVarP(&address, "host", "H", "", "Address to run server on - default value is from viper config")
	serverCommand.Flags().UintVarP(&port, "port", "P", 0, "Port to run server on - default value is from viper config")
	clientCommand.Flags().StringP("host", "H", "", "Host to send form to.")
	clientCommand.Flags().UintP("port", "P", 0, "Port to send form to.")
	rootCmd.AddCommand(clientCommand, serverCommand, versionCommand)
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}

}
