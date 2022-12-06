package main

import (
	_ "embed"
	"os"

	"github.com/dapr/go-sdk/client"
	daprd "github.com/dapr/go-sdk/service/http"
	"github.com/spolab/petclinic/owner/server"
	"go.uber.org/zap"
)

//go:embed version.txt
var revision string

func main() {
	stateStoreName := os.Getenv("STATESTORE_NAME")
	pubsubName := os.Getenv("PUBSUB_NAME")
	pubsubTopic := os.Getenv("PUBSUB_TOPIC")
	//
	// Initialize logging
	//
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	//
	// Announces the bootstrap of the microservice
	//
	logger.Info("Starting Owner Microservice", zap.String("revision", revision), zap.String("statestore_name", stateStoreName), zap.String("pubsub_name", pubsubName), zap.String("pubsub_topic", pubsubTopic))
	//
	// Initialize the DAPR Client
	//
	dapr, err := client.NewClient()
	if err != nil {
		logger.Panic("Error initializing the Dapr client", zap.Error(err))
	}
	defer dapr.Close()
	//
	// Initialize the command actor factory
	//
	service := daprd.NewService("127.0.0.1:3000")
	service.RegisterActorImplFactory(server.OwnerActorFactory(logger, dapr))
	//
	// Start the server
	//
	if err := service.Start(); err != nil {
		logger.Panic("starting the http service", zap.Error(err))
	}
	service.GracefulStop()
}
