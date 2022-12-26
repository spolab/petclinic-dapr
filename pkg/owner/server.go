package owner

import (
	"context"

	"github.com/dapr/go-sdk/actor"
	"github.com/dapr/go-sdk/client"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	"github.com/spolab/petstore/gen/api"
	"github.com/spolab/petstore/pkg/event"
)

type OwnerState struct {
	Salutation string `json:"salutation"`
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Address    string `json:"address"`
	PostCode   string `json:"post_code"`
	City       string `json:"city"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
}

type OwnerActor struct {
	actor.ServerImplBase
	dapr     client.Client
	validate *validator.Validate
	broker   string
	topic    string
}

// OwnerActorFactory is a
func OwnerActorFactory(dapr client.Client, broker string, topic string) actor.Factory {
	return func() actor.Server {
		log.Info().Interface("client", dapr).Str("broker", broker).Str("topic", topic).Msg("initializing actor server")
		return &OwnerActor{dapr: dapr, validate: validator.New(), broker: broker, topic: topic}
	}
}

// Type returns the name that will be used to bind the actor to the URL, i.e. `/actor/<type>/method/<methodName>`
func (a *OwnerActor) Type() string {
	return "owner"
}

const keyContactState = "contact"

func (a *OwnerActor) Register(ctx context.Context, cmd *api.RegisterOwnerCommand) *api.RegisterOwnerResponse {
	log.Info().Str("id", a.ID()).Msg("start registering new owner")
	// validate the command
	if err := a.validate.Struct(cmd); err != nil {
		return &api.RegisterOwnerResponse{Result: api.CommandResult_INVALID, Message: err.Error()}
	}
	// exit if the owner already exists
	found, err := a.GetStateManager().Contains(keyContactState)
	if err != nil {
		return &api.RegisterOwnerResponse{Result: api.CommandResult_ERROR, Message: err.Error()}
	}
	if found {
		return &api.RegisterOwnerResponse{Result: api.CommandResult_INVALID, Message: "owner already exists"}
	}
	// store the state of the owner
	state := OwnerState{
		Salutation: cmd.Salutation.String(),
		Name:       cmd.Name,
		Surname:    cmd.Surname,
	}
	if err := a.GetStateManager().Set(keyContactState, &state); err != nil {
		return &api.RegisterOwnerResponse{Result: api.CommandResult_ERROR, Message: err.Error()}
	}
	// all good, publish the event
	if err := a.dapr.PublishEvent(ctx, a.broker, a.topic, event.OwnerRegistered{
		ID:         a.ID(),
		Salutation: cmd.Salutation.String(),
		Name:       cmd.Name,
		Surname:    cmd.Surname,
	}); err != nil {
		return &api.RegisterOwnerResponse{Result: api.CommandResult_ERROR, Message: err.Error()}
	}
	return &api.RegisterOwnerResponse{Result: api.CommandResult_OK, Message: ""}
}
