package app

import (
	"fmt"
	"net/http"
	"os"

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
func (a *App) Start() error {
	log.Info().Msgf("Server is listening on http://%s:%d", a.Server.Host, a.Server.Port)

	a.registerHandlers()

	err := http.ListenAndServe(fmt.Sprintf("%s:%d", a.Server.Host, a.Server.Port), nil)
	if err != nil {
		return err
	}

	return nil
}

// registerHandlers registers all handlers with their respective paths in the HTTP router of this application's
func (a *App) registerHandlers() {
	http.HandleFunc("/api/v1/users/", a.createUserHandler)
}
