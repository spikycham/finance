package network

import (
	"encoding/json"
	"net/http"
)

type StandardResponse[T any] struct {
	Message *string
	Data    *T
	Error   *string
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
