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
package vet

import (
	"context"
	"fmt"

	"github.com/dapr/go-sdk/actor"
	"github.com/dapr/go-sdk/client"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	"github.com/spolab/petstore/pkg/common/events"
)

const keyDetails = "details"

type ActorDetails struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
	Active  bool   `json:"active"`
}

type Actor struct {
	actor.ServerImplBase
	dapr     client.Client
	validate *validator.Validate
	broker   string
	topic    string
}

func (*Actor) Type() string {
	return "vet"
}

func ActorFactory(dapr client.Client, validator *validator.Validate, broker string, topic string) actor.Factory {
	return func() actor.Server {
		log.Info().Msg("activating actor")
		return &Actor{dapr: dapr, validate: validator, broker: broker, topic: topic}
	}
}

func (actor *Actor) Register(ctx context.Context, cmd *RegisterVetCommand) (*RegisterVetResponse, error) {
	log.Info().Str("id", actor.ID()).Msg("begin register")
	//
	// Return an error if the command does not pass validation
	//
	log.Debug().Str("id", actor.ID()).Msg("validate command")
	if err := actor.validate.Struct(cmd); err != nil {
		log.Error().Str("id", actor.ID()).Err(err).Msg("validating command")
		return nil, err
	}
	//
	// Return an error if the actor instance already exists
	//
	log.Debug().Str("id", actor.ID()).Msg("check if vet already exists")
	found, err := actor.GetStateManager().Contains(keyDetails)
	if err != nil {
		log.Error().Str("id", actor.ID()).Err(err).Msg("executing statemanager::contains")
		return nil, err
	}
	if found {
		log.Error().Str("id", actor.ID()).Msg("vet already exists")
		return nil, fmt.Errorf("vet '%s' already exists", actor.ID())
	}
	//
	// Stores the state of the aggregate
	//
	log.Debug().Str("id", actor.ID()).Msg("snapshot aggregate state")
	details := &ActorDetails{Id: actor.ID(), Name: cmd.Name, Surname: cmd.Surname, Phone: cmd.Phone, Email: cmd.Email, Active: true}
	if err := actor.GetStateManager().Set(keyDetails, details); err != nil {
		log.Error().Str("id", actor.ID()).Err(err).Msg("snapshotting aggregate state")
		return nil, err
	}
	//
	// Now that we know that the state is safely stored, let´s broadcast the event
	// NOTE: I know this is not fail-safe. This code is just for illustrative purposes. Will be improved later.
	//
	log.Debug().Str("id", actor.ID()).Str("broker", actor.broker).Str("topic", actor.topic).Msg("publish event")
	event := events.CloudEvent("vet", "VetRegistered", &VetRegistered{Id: actor.ID(), Name: cmd.Name, Surname: cmd.Surname, Phone: cmd.Phone, Email: cmd.Email})
	if err := actor.dapr.PublishEvent(ctx, actor.broker, actor.topic, event); err != nil {
		log.Error().Str("id", actor.ID()).Err(err).Msg("publishing event")
		return nil, err
	}
	log.Info().Str("id", actor.ID()).Msg("end register")
	return &RegisterVetResponse{Status: 0, Message: "OK"}, nil
}
