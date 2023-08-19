package main

import (
	"flag"
	"fmt"
	"os"

	envserver "github.com/Mahmoud-Emad/envserver/app"
)

func main() {
	var configFilePath string
	flag.StringVar(&configFilePath, "config", "", "Path to the configuration file")
	flag.Parse()

	if configFilePath == "" {
		fmt.Println("Error: You must provide the path to the configuration file using the -config flag.")
		flag.Usage()
		os.Exit(1)
	}

	fmt.Println("configFilePath: ", configFilePath)
	app, err := envserver.NewApp(configFilePath)
	if err != nil {
		fmt.Printf("Error creating the app: %s\n", err)
		os.Exit(1)
	}

	app.Start()
}
