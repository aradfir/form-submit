package main

import (
	"FormSubmit/client"
	"FormSubmit/server"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Enter server,client or version!")
		os.Exit(0)
	}
	functionality := os.Args[1]
	switch functionality {
	case "server":
		server.RunServer()
	case "client":
		client.RunClient()
	case "version":
		fmt.Println("version 1")
	default:
		fmt.Println("Invalid args!")
	}
}
