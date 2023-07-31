package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	internal "github.com/Mahmoud-Emad/envserver/internal"
	models "github.com/Mahmoud-Emad/envserver/models"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// App for all dependencies of backend server
type App struct {
	Config internal.Configuration
	Server Server
	DB     internal.Database
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
	db := internal.NewDatabase()
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

// Define a custom type for the context key to avoid potential collisions with other keys.
type contextKey int

// Create a new context key to store the user information.
const UserContextKey contextKey = 1

// Get the requested user data.
func (a *App) GetRequestedUser(r *http.Request) (models.User, error) {
	user, ok := r.Context().Value(UserContextKey).(models.User)
	if !ok {
		authHeader := r.Header.Get("Authorization")
		user, err := VerifyAndDecodeJwtToken(authHeader, a.Config.Server.JWTSecretKey)
		if err != nil {
			return user, errors.New("Cannot decode jwt.")
		}

		// Add the user object inside the request.
		ctx := context.WithValue(r.Context(), UserContextKey, user)
		r = r.WithContext(ctx)
	}
	return user, nil
}

func (a *App) registerHandlers() {
	r := mux.NewRouter()

	// API version 1 routes
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
