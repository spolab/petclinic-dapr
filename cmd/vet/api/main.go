package main

import (
	_ "embed"
	"os"

	"github.com/dapr/go-sdk/client"
	"github.com/dapr/go-sdk/service/http"
	"github.com/rs/zerolog/log"
)

var revision string

func main() {
	pubsub := os.Getenv("PUBSUB_NAME")
	topic := os.Getenv("PUBSUB_TOPIC")
	//
	// Announces the bootstrap of the microservice
	//
	log.Info().Str("revision", revision).Str("pubsub", pubsub).Str("topic", topic).Msg("starting owner microservice")
	//
	// Connect to the DAPR sidecar
	//
	_, err := client.NewClient()
	if err != nil {
		log.Fatal().Err(err).Msg("connecting to the dapr sidecar")
	}
	//
	// Initialize the router
	//
	app := http.NewService("127.0.0.1:3000")
	//
	// Start the server
	//
	log.Fatal().Err(app.Start()).Msg("starting the http service")
}
