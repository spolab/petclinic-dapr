/*
Copyright 2022 Alessandro Santini

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
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
