package cmd

import (
	"FormSubmit/internal/config"
	"FormSubmit/internal/server"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"net/http"
)

var (
	address string
	port    uint
)
var serverCommand = &cobra.Command{
	Use:   "server",
	Short: "Runs the server",
	Args:  cobra.NoArgs,

	Run: serve,
}

func init() {
	serverCommand.Flags().StringVarP(&address, "host", "H", "", "Address to run server on - default value is from viper config")
	serverCommand.Flags().UintVarP(&port, "port", "P", 0, "Port to run server on - default value is from viper config")
	RootCmd.AddCommand(serverCommand)
}

func SetupMetricServer() {
	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(":2112", nil)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Error starting metrics server")
	}
}

func serve(cmd *cobra.Command, args []string) {
	go SetupMetricServer()
	cfg, err := config.ViperSetup()
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Error reading config")
	}

	address, port = config.GetHostAndPort(address, port, &cfg)
	server.RunServer(address, port)
}
