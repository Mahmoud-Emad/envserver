package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	internal "github.com/Mahmoud-Emad/envserver/internal"
	models "github.com/Mahmoud-Emad/envserver/models"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// App for all dependencies of backend server
type App struct {
	Config internal.Configuration
	Server Server
	DB     models.Database
}

func initZerolog() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

// NewApp creates a new App instance using the provided configuration file.
func NewApp(configFileName string) (app *App, err error) {
	initZerolog()

	config, err := internal.ReadConfigFromFile(configFileName)
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, err
	}

	server := NewServer(config.Server.Host, config.Server.Port)
	db := models.NewDatabase()
	err = db.Connect(config.Database)

	if err != nil {
		return nil, err
	}

	err = db.Migrate()
	if err != nil {
		return
	}

	return &App{
		Server: *server,
		Config: config,
		DB:     db,
	}, nil
}

// Start starts the server and listens for incoming requests.
func (a *App) Start() {
	// Create a channel to listen for OS signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	a.registerHandlers()

	// Create a new server
	server := &http.Server{
		Addr: string(rune(a.Server.Port)),
	}
	// Start the server in a goroutine
	go func() {
		log.Info().Msgf("Server is listening on http://%s:%d", a.Server.Host, a.Server.Port)
		if err := http.ListenAndServe(fmt.Sprintf("%s:%d", a.Server.Host, a.Server.Port), nil); err != nil && err != http.ErrServerClosed {
			log.Error().Msgf("Server failed to start: %v", err)
		}
		log.Info().Msg("Stopped serving new connections")
	}()

	// Wait for the shutdown signal
	<-stop

	// Create a context with a timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown the server
	if err := server.Shutdown(ctx); err != nil {
		log.Error().Msgf("Server shutdown error: %v", err)
	} else {
		log.Info().Msgf("Server gracefully stopped")
	}
}

// registerHandlers registers all handlers with their respective paths in the HTTP router of this application's
func (a *App) registerHandlers() {
	http.HandleFunc("/api/v1/users/", a.createUserHandler)
}
