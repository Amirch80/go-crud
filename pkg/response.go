package pkg

import (
	"encoding/json"
	"net/http"
)

func ResponseJson(writer http.ResponseWriter, status int, data any) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	if data != nil {
		_ = json.NewEncoder(writer).Encode(data)
	}
}

func ResponseError(writer http.ResponseWriter, status int, message string) {
	ResponseJson(writer, status, map[string]string{
		"error": message,
	})
}
