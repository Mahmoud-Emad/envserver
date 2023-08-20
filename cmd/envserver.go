package main

import (
	"flag"
	"os"

	envserver "github.com/Mahmoud-Emad/envserver/app"
	"github.com/rs/zerolog/log"
)

func main() {
	var configFilePath string
	flag.StringVar(&configFilePath, "config", "", "Path to the Config file")
	flag.Parse()

	if configFilePath == "" {
		log.Error().Msgf("Error: You must provide the path to the Config file using the -config flag.")
		flag.Usage()
		os.Exit(1)
	}

	app, err := envserver.NewApp(configFilePath)
	if err != nil {
		log.Error().Msgf("Error creating the app: %s\n", err)
		os.Exit(1)
	}

	app.Start()
}
