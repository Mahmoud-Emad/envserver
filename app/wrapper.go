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
