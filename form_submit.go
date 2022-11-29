package main

import (
	"FormSubmit/cmd"
	"fmt"
	"os"
)

func main() {

	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}

}
