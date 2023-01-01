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

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/dapr/go-sdk/actor"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	"github.com/spolab/petstore/pkg/command"
	"github.com/spolab/petstore/pkg/event"
)

const (
	KeyState             = "state"
	KeyEvents            = "events"
	EventVetRegisteredV1 = "VetRegistered/v1"
)

type VetState struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
	Deleted bool   `json:"deleted"`
	Version int    `json:"version"`
}

type Vet struct {
	actor.ServerImplBase
	state    *VetState
	validate *validator.Validate
	events   []*cloudevents.Event
}

func (*Vet) Type() string {
	return "vet"
}

func VetActorFactory(v *validator.Validate) actor.Factory {
	return func() actor.Server {
		log.Info().Msg("activating actor")
		return &Vet{validate: v}
	}
}

// Registers a new vet. Bear in mind that the method sigmature must contain the error, even if it is currently mishandled (v1.9.5)
func (vet *Vet) Register(ctx context.Context, cmd *command.RegisterVetCommand) (*command.ActorResponse, error) {
	log.Info().Str("id", vet.ID()).Msg("begin register")
	if err := vet.Load(); err != nil {
		return &command.ActorResponse{Status: command.StatusError, Message: err.Error()}, nil
	}
	//
	// Trying to create a vet that already has events is clearly an error, no matter how valid the command is. Kill it.
	//
	if len(vet.events) > 0 {
		return &command.ActorResponse{Status: command.StatusInvalid, Message: "vet already exists"}, nil
	}
	//
	// Return an error if the command does not pass validation
	//
	log.Debug().Str("id", vet.ID()).Msg("validate command")
	if err := vet.validate.Struct(cmd); err != nil {
		log.Error().Str("id", vet.ID()).Err(err).Msg("validating command")
		return &command.ActorResponse{Status: command.StatusInvalid, Message: err.Error()}, nil
	}
	//
	// Now that we know that the state is safely stored, let´s broadcast the event
	// NOTE: I know this is not fail-safe. This code is just for illustrative purposes. Will be improved in another edition.
	//
	log.Debug().Str("id", vet.ID()).Msg("store actor events")
	event := event.CloudEvent("vet", EventVetRegisteredV1, &event.VetRegistered{
		Id:      vet.ID(),
		Name:    cmd.Name,
		Surname: cmd.Surname,
		Phone:   cmd.Phone,
		Email:   cmd.Email})
	//
	// Apply the event to alter the state
	//
	if err := vet.Apply(&event); err != nil {
		return &command.ActorResponse{Status: command.StatusError, Message: err.Error()}, nil
	}
	//
	// Append the event to the queue of events to be committed
	//
	if err := vet.Append(&event); err != nil {
		log.Error().Str("id", vet.ID()).Err(err).Msg("storing actor events")
		return &command.ActorResponse{Status: command.StatusError, Message: err.Error()}, nil
	}
	//
	// Return the events as response
	//
	log.Info().Str("id", vet.ID()).Msg("end register")
	return &command.ActorResponse{Status: command.StatusOK, Events: []*cloudevents.Event{&event}}, nil
}

// Apply alters the state of the aggregate based on the event
func (vet *Vet) Apply(src *cloudevents.Event) error {
	switch src.Type() {
	case EventVetRegisteredV1:
		ev := event.VetRegistered{}
		if err := src.DataAs(&ev); err != nil {
			return err
		}
		vet.state.Id = ev.Id
		vet.state.Surname = ev.Surname
		vet.state.Name = ev.Name
		vet.state.Phone = ev.Phone
		vet.state.Email = ev.Email
	default:
		return fmt.Errorf("unknown event type '%s'", src.Type())
	}
	return nil
}

// Load the state of the aggregate from the events log. Returns
func (vet *Vet) Load() error {
	if vet.state == nil {
		vet.state = &VetState{}
		err := vet.GetStateManager().Get(KeyEvents, &vet.events)
		if err != nil {
			return err
		}
		for index, event := range vet.events {
			vet.state.Version = index
			if err := vet.Apply(event); err != nil {
				//
				// Erase the state if an error occurs.
				// In this way, any further attempt at using this actor instance will fail until the event stream handling gets fixed.
				//
				vet.state = nil
				return err
			}
		}
	}
	return nil
}

func (vet *Vet) Append(events ...*cloudevents.Event) error {
	vet.events = append(vet.events, events...)
	return vet.GetStateManager().Set(KeyEvents, vet.events)
}
