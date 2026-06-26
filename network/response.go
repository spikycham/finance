package network

import (
	"encoding/json"
	"errors"
	"net/http"
)

type StandardResponse[T any] struct {
	Message *string `json:"message"`
	Data    *T      `json:"data"`
	Error   *string `json:"error"`
}

var (
	ErrInternal      = errors.New("internal error")
	ErrInvalidJSON   = errors.New("invalid JSON body")
	ErrMissingFields = errors.New("missing fields")
)

func ReadBody[T any](r *http.Request) (*T, error) {
	defer r.Body.Close()

	var body T
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, ErrInternal
	}

	return &body, nil
}

func ResponseMessage(w http.ResponseWriter, code int, msg string) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(&StandardResponse[any]{
		Message: &msg,
		Error:   nil,
		Data:    nil,
	})
}

func ResponseError(w http.ResponseWriter, code int, err string) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(&StandardResponse[any]{
		Message: nil,
		Error:   &err,
		Data:    nil,
	})
}

func ResponseJSON[T any](w http.ResponseWriter, code int, data T) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(&StandardResponse[T]{
		Message: nil,
		Error:   nil,
		Data:    &data,
	})
}
