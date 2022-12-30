package main

import (
	_ "embed"
	"net/http"
	"os"
	"time"

	"github.com/dapr/go-sdk/client"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"github.com/spolab/petstore/pkg/vet"
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
	router.Put("/{id}", vet.Register(dapr))
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
