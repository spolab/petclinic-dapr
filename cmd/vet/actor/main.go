package main

import (
	"github.com/dapr/go-sdk/client"
	"github.com/dapr/go-sdk/service/http"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	"github.com/spolab/petstore/pkg/vet"
)

func main() {
	client, err := client.NewClient()
	if err != nil {
		log.Fatal().Err(err).Msg("connecting to the dapr sidecar")
	}
	app := http.NewService("127.0.0.1:3000")
	app.RegisterActorImplFactory(vet.ActorFactory(client, validator.New()))
	log.Fatal().Err(app.Start()).Msg("starting the actor server")
}
