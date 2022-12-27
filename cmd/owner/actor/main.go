package main

import (
	_ "embed"
	"os"

	"github.com/dapr/go-sdk/client"
	daprd "github.com/dapr/go-sdk/service/http"
	"github.com/rs/zerolog/log"
	"github.com/spolab/petstore/pkg/owner"
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
	// Initialize the DAPR Client
	//
	dapr, err := client.NewClient()
	if err != nil {
		log.Fatal().Err(err).Msg("Error initializing the Dapr client")
	}
	defer dapr.Close()
	//
	// Initialize the command actor factory
	//
	service := daprd.NewService("127.0.0.1:3000")
	defer service.GracefulStop()
	service.RegisterActorImplFactory(owner.OwnerActorFactory(dapr, pubsub, topic))
	//
	// Start the server
	//
	log.Fatal().Err(service.Start()).Msg("starting the http service")
}
