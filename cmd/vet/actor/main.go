package main

import (
	"os"

	"github.com/dapr/go-sdk/client"
	"github.com/dapr/go-sdk/service/http"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	"github.com/spolab/petstore/pkg/vet"
)

func main() {
	//
	// Load the start parameters
	//
	broker := os.Getenv("BROKER")
	topic := os.Getenv("TOPIC")
	log.Info().Str("pubsub", broker).Str("topic", topic).Msg("starting owner microservice")
	//
	// Connect to the DAPR sidecar
	//
	client, err := client.NewClient()
	if err != nil {
		log.Fatal().Err(err).Msg("connecting to the dapr sidecar")
	}
	//
	// Start the actor server
	//
	app := http.NewService("127.0.0.1:3000")
	app.RegisterActorImplFactory(vet.ActorFactory(client, validator.New(), broker, topic))
	log.Fatal().Err(app.Start()).Msg("starting the actor server")
}
