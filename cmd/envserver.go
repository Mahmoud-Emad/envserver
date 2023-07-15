// // This file just for testing.
package main

import (
	envserver "github.com/Mahmoud-Emad/envserver/app"
)

func main() {
	configFileName := "config.toml"
	app, err := envserver.NewApp(configFileName)
	if err != nil {
		return
	}
	app.Start()
}
