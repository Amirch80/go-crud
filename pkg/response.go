package pkg

import (
	"crud-task/pkg/appError"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type Response struct {
	Message string            `json:"message,omitempty"`
	Data    any               `json:"data,omitempty"`
	Errors  map[string]string `json:"errors,omitempty"`
}

func ResponseJson(writer http.ResponseWriter, status int, response *Response) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	if response != nil {
		_ = json.NewEncoder(writer).Encode(response)
	}
}

func ResponseError(writer http.ResponseWriter, status int, message string) {
	ResponseJson(writer, status, &Response{
		Message: message,
	})
}

func ResponseValidationError(writer http.ResponseWriter, status int, errors map[string]string) {
	ResponseJson(writer, status, &Response{
		Message: "validation failed",
		Errors:  errors,
	})
}

func ResponseFromError(writer http.ResponseWriter, error error) {
	var appErr *appError.AppError
	if errors.As(error, &appErr) {
		if appErr.Status >= 500 {
			log.Printf("internal error: %v", appErr.Err)
		}
		ResponseError(writer, appErr.Status, appErr.Message)
		return
	}

	log.Printf("unhandled error: %v", error)
	ResponseError(writer, http.StatusInternalServerError, "internal server error")
}
