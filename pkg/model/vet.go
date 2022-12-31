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
package model

import (
	"context"
	"fmt"

	"github.com/dapr/go-sdk/actor"
	"github.com/dapr/go-sdk/client"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	"github.com/spolab/petstore/pkg/command"
	"github.com/spolab/petstore/pkg/event"
)

const keyDetails = "details"

type VetDetails struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
	Active  bool   `json:"active"`
}

type Vet struct {
	actor.ServerImplBase
	dapr     client.Client
	validate *validator.Validate
	broker   string
	topic    string
}

func (*Vet) Type() string {
	return "vet"
}

func VetActorFactory(dapr client.Client, validator *validator.Validate, broker string, topic string) actor.Factory {
	return func() actor.Server {
		log.Info().Msg("activating actor")
		return &Vet{dapr: dapr, validate: validator, broker: broker, topic: topic}
	}
}

func (vet *Vet) Register(ctx context.Context, cmd *command.RegisterVetCommand) (*command.RegisterVetResponse, error) {
	log.Info().Str("id", vet.ID()).Msg("begin register")
	//
	// Return an error if the command does not pass validation
	//
	log.Debug().Str("id", vet.ID()).Msg("validate command")
	if err := vet.validate.Struct(cmd); err != nil {
		log.Error().Str("id", vet.ID()).Err(err).Msg("validating command")
		return &command.RegisterVetResponse{Status: command.StatusInvalid, Message: err.Error()}, nil
	}
	//
	// Return an error if the actor instance already exists
	//
	log.Debug().Str("id", vet.ID()).Msg("check if vet already exists")
	found, err := vet.GetStateManager().Contains(keyDetails)
	if err != nil {
		log.Error().Str("id", vet.ID()).Err(err).Msg("executing statemanager::contains")
		return &command.RegisterVetResponse{Status: command.StatusError, Message: err.Error()}, nil
	}
	if found {
		log.Error().Str("id", vet.ID()).Msg("vet already exists")
		return &command.RegisterVetResponse{Status: command.StatusInvalid, Message: fmt.Sprintf("vet '%s' already exists", vet.ID())}, nil
	}
	//
	// Stores the state of the aggregate
	//
	log.Debug().Str("id", vet.ID()).Msg("snapshot aggregate state")
	details := &VetDetails{Id: vet.ID(), Name: cmd.Name, Surname: cmd.Surname, Phone: cmd.Phone, Email: cmd.Email, Active: true}
	if err := vet.GetStateManager().Set(keyDetails, details); err != nil {
		log.Error().Str("id", vet.ID()).Err(err).Msg("snapshotting aggregate state")
		return &command.RegisterVetResponse{Status: command.StatusError, Message: err.Error()}, nil
	}
	//
	// Now that we know that the state is safely stored, let´s broadcast the event
	// NOTE: I know this is not fail-safe. This code is just for illustrative purposes. Will be improved later.
	//
	log.Debug().Str("id", vet.ID()).Str("broker", vet.broker).Str("topic", vet.topic).Msg("publish event")
	event := event.CloudEvent("vet", "VetRegistered", &event.VetRegistered{Id: vet.ID(), Name: cmd.Name, Surname: cmd.Surname, Phone: cmd.Phone, Email: cmd.Email})
	if err := vet.dapr.PublishEvent(ctx, vet.broker, vet.topic, event); err != nil {
		log.Error().Str("id", vet.ID()).Err(err).Msg("publishing event")
		return &command.RegisterVetResponse{Status: command.StatusError, Message: err.Error()}, nil
	}
	log.Info().Str("id", vet.ID()).Msg("end register")
	return &command.RegisterVetResponse{Status: command.StatusOK, Message: "OK"}, nil
}
