package respond

import (
	"encoding/json"
	"net/http"
)

const (
	keyContentType  = "Content-Type"
	applicationJSON = "application/json"
	textPlain       = "text/plain"
)

func JSON(w http.ResponseWriter, statusCode int, source any) error {
	bytes, err := json.Marshal(source)
	if err != nil {
		return err
	}
	return respond(w, statusCode, applicationJSON, bytes)
}

func String(w http.ResponseWriter, statusCode int, source string) error {
	return respond(w, statusCode, textPlain, []byte(source))
}

func respond(w http.ResponseWriter, statusCode int, contentType string, data []byte) error {
	w.WriteHeader(statusCode)
	w.Header().Add(keyContentType, contentType)
	if _, err := w.Write(data); err != nil {
		return err
	}
	return nil
}
