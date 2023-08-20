package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	internal "github.com/Mahmoud-Emad/envserver/internal"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// App for all dependencies of backend server
type App struct {
	Config internal.Config
	Server Server
	DB     internal.Database
}

func initZerolog() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

// NewApp creates a new App instance using the provided Config file.
func NewApp(configFileName string) (*App, error) {
	initZerolog()

	log.Info().Msg("Loading config file.")
	config, err := internal.ReadConfigFromFile(configFileName)
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, err
	}
	log.Info().Msg("Config file loaded.")

	server := NewServer(config.Server.Host, config.Server.Port)
	db := internal.NewDatabase()
	err = db.Connect(config.Database)

	if err != nil {
		return nil, err
	}

	err = db.Migrate()
	if err != nil {
		return nil, err
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
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(a.Config.Server.ShutdownTimeout)*time.Second)
	defer cancel()

	// Shutdown the server
	if err := server.Shutdown(ctx); err != nil {
		log.Error().Msgf("Server shutdown error: %v", err)
	}
	log.Info().Msgf("Server gracefully stopped")
}

func (a *App) registerHandlers() {
	r := mux.NewRouter()

	apiRouter := r.PathPrefix("/api/v1").Subrouter()
	userRouter := apiRouter.PathPrefix("/users").Subrouter()
	authRouter := apiRouter.PathPrefix("/auth").Subrouter()
	projectRouter := apiRouter.PathPrefix("/projects").Subrouter()

	// User routes (protected with authentication)
	userRouter.HandleFunc("", a.wrapRequest(a.getUsersHandler, true)).Methods(http.MethodGet, http.MethodOptions)
	userRouter.HandleFunc("/{id}", a.wrapRequest(a.deleteUserByIDHandler, true)).Methods(http.MethodDelete, http.MethodOptions)

	// Auth routes
	authRouter.HandleFunc("/signup", a.wrapRequest(a.signupHandler, false)).Methods(http.MethodPost, http.MethodOptions)
	authRouter.HandleFunc("/signin", a.wrapRequest(a.signinHandler, false)).Methods(http.MethodPost, http.MethodOptions)

	// Project routes (protected with authentication)
	projectRouter.HandleFunc("", a.wrapRequest(a.createProjectHandler, true)).Methods(http.MethodPost, http.MethodOptions)
	projectRouter.HandleFunc("", a.wrapRequest(a.getProjectsHandler, true)).Methods(http.MethodGet, http.MethodOptions)
	projectRouter.HandleFunc("/{id}", a.wrapRequest(a.getProjectByIDHandler, true)).Methods(http.MethodGet, http.MethodOptions)
	projectRouter.HandleFunc("/{id}", a.wrapRequest(a.deleteProjectByIDHandler, true)).Methods(http.MethodDelete, http.MethodOptions)

	// Add the authentication middleware to the protected routes
	userRouter.Use(a.authenticateMiddleware)
	projectRouter.Use(a.authenticateMiddleware)

	// Set the router for the application
	http.Handle("/", r)
}
