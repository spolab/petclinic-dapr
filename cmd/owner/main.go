package main

import (
	_ "embed"
	"os"

	"github.com/gofiber/fiber/v2"
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
	// Initialize the router
	//
	app := fiber.New()
	//
	// Start the server
	//
	log.Fatal().Err(app.Listen(":3000")).Msg("starting the http service")
}
