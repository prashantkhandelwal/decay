package main

import (
	"flag"

	"github.com/prashantkhandelwal/decay/config"
	"github.com/prashantkhandelwal/decay/server"
)

func main() {

	port := flag.String("port", "8989", "Specify the port to run the server.")
	env := flag.String("env", "release", "Switch between release or debug mode.")

	flag.Parse()

	config := config.Config{
		Environment: *env,
		Port:        *port,
	}

	// Starts the server
	server.Run(&config)
}
