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

type VetRegistered struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
}

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

func (Actor) Type() string {
	return "vet"
}

func ActorFactory(dapr client.Client, validator *validator.Validate, broker string, topic string) actor.Factory {
	return func() actor.Server {
		log.Info().Msg("activating actor")
		return &Actor{dapr: dapr, validate: validator}
	}
}

func (actor *Actor) Register(ctx context.Context, cmd *RegisterVetCommand) error {
	//
	// Return an error if the command does not pass validation
	//
	if err := actor.validate.Struct(cmd); err != nil {
		return err
	}
	//
	// Return an error if the actor instance already exists
	//
	found, err := actor.GetStateManager().Contains(keyDetails)
	if err != nil {
		return err
	}
	if found {
		return fmt.Errorf("vet '%s' already exists", actor.ID())
	}
	//
	// Stores the state of the aggregate
	//
	details := &ActorDetails{Id: actor.ID(), Name: cmd.Name, Surname: cmd.Surname, Phone: cmd.Phone, Email: cmd.Email, Active: true}
	if err := actor.GetStateManager().Set(keyDetails, details); err != nil {
		return err
	}
	//
	// Now that we know that the state is safely stored, let´s broadcast the event
	// NOTE: I know this is not fail-safe. This code is just for illustrative purposes. Will be improved later.
	//
	event := events.CloudEvent("vet", "VetRegistered", &VetRegistered{Id: actor.ID(), Name: cmd.Name, Surname: cmd.Surname, Phone: cmd.Phone, Email: cmd.Email})
	if err := actor.dapr.PublishEvent(ctx, actor.broker, actor.topic, event); err != nil {
		return err
	}
	return nil
}
