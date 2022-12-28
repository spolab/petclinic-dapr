package vet

import (
	"context"
	"fmt"

	"github.com/dapr/go-sdk/actor"
	"github.com/dapr/go-sdk/client"
	"github.com/go-playground/validator/v10"
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
}

func (Actor) Type() string {
	return "vet"
}

func ActorFactory(dapr client.Client, validator *validator.Validate) actor.Factory {
	return func() actor.Server {
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
	actor.GetStateManager().Set(keyDetails, details)
	return nil
}
