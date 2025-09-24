package main

import (
	"fmt"
	"log"
	"os"

	"github.com/alecthomas/kingpin/v2"
	"github.com/prashantkhandelwal/decay/config"
	"github.com/prashantkhandelwal/decay/server"
)

func main() {

	app := kingpin.New("decay", "Share files with a decay timer.")
	port := app.Flag("port", "Specify the port to run the server.").Default("8989").String()
	//env := kingpin.Flag("env", "Switch between release or debug mode.").Default("release").String()
	//appConfig := app.Flag("config", "Path to config file.").Default("config.yaml").String()

	//configInit := kingpin.Command("init", "Initialize the config file.")

	_, err := app.Parse(os.Args[1:])
	if err != nil {
		log.Fatalf("ERROR: Failed to parse flags: %v", err)
	}

	app.HelpFlag.Short('h')
	app.Version("decay 0.1.0")

	c, err := config.InitConfig()
	if err != nil {
		log.Fatalf("ERROR: Cannot load configuration = %v", err)
	}

	fmt.Printf("Using config file: %v\n", c)

	if *port != "8989" {
		c.Server.PORT = *port
	}

	// Starts the server
	server.Run(c)
}
