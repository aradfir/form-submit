package cmd

import "github.com/spf13/cobra"

var RootCmd = &cobra.Command{
	Use:   "FormSubmit",
	Short: "Form submit is a sample project to work with go, cobra, viper and grpc",
	Run: func(cmd *cobra.Command, args []string) {

	},
}
