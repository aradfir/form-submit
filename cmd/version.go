package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var commit string
var buildTime string
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Get the version of this program",
	Long:  "Info about the program aka build sha and build time",
	Run:   printVersion,
}

func init() {
	RootCmd.AddCommand(versionCmd)
}

func printVersion(cmd *cobra.Command, args []string) {
	fmt.Println("Commit SHA:", commit)
	fmt.Println("Build time:", buildTime)
}
