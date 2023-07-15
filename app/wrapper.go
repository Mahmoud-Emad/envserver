package app

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
)

type Response struct {
	Message string        `json:"message"`
	Status  int           `json:"status"`
	Data    interface{}   `json:"data"`
	Error   *ErrorDetails `json:"error,omitempty"`
}

type ErrorDetails struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ResponeError error
type Handler func(w http.ResponseWriter, r *http.Request)

func sendJSONResponse(w http.ResponseWriter, status int, message string, data interface{}, err error) {

	errDetails := &ErrorDetails{}
	response := Response{
		Status:  status,
		Message: message,
		Data:    data,
	}

	if err != nil {
		log.Error().Msg(fmt.Sprintf("%v", err))
		errDetails.Code = status
		errDetails.Message = err.Error()
		response.Error = errDetails
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}

// Used for debugging, print out the incoming request URL and its method.
func wrapRequest(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Debug().Msgf("Request: %s %s", r.Method, r.URL.Path)
		h(w, r)
	}
}
