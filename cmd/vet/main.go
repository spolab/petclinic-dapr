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
package main

import (
	_ "embed"
	"net/http"
	"os"
	"time"

	"github.com/dapr/go-sdk/client"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"github.com/spolab/petstore/pkg/api"
)

var revision string

func main() {
	pubsub := os.Getenv("BROKER")
	topic := os.Getenv("TOPIC")
	//
	// Announces the bootstrap of the microservice
	//
	log.Info().Str("revision", revision).Str("pubsub", pubsub).Str("topic", topic).Msg("starting owner microservice")
	//
	// Connect to the DAPR sidecar
	//
	dapr, err := client.NewClient()
	if err != nil {
		log.Fatal().Err(err).Msg("connecting to the dapr sidecar")
	}
	//
	// Initialize the router
	//
	router := chi.NewRouter()
	router.Put("/{id}", api.Register(dapr))
	//
	// Start the server
	//
	app := http.Server{
		Addr:         ":3000",
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Fatal().Err(app.ListenAndServe()).Msg("starting the http service")
}
