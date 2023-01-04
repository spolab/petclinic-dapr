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
	"github.com/dapr/go-sdk/service/http"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	"github.com/spolab/petstore/pkg/model"
)

func main() {
	log.Info().Msg("starting vet microservice")
	//
	// Start the actor server
	//
	app := http.NewService("127.0.0.1:3000")
	app.RegisterActorImplFactory(model.VetActorFactory(validator.New()))
	app.RegisterActorImplFactory(model.ClientActorFactory)
	log.Fatal().Err(app.Start()).Msg("starting the actor server")
}
