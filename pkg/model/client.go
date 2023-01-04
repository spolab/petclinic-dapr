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

type ClientState struct {
	Salutation string
	Surname    string
	Name       string
	Phone      string
	Email      string
	Version    int
}

type Client struct {
	BaseAggregateRoot[ClientState]
	validate *validator.Validate
}

func (actor *Client) Type() string {
	return "client"
}

// register a new client
func (actor *Client) Register(ctx context.Context, cmd *command.RegisterClientCommand) (*command.ActorResponse, error) {
	res, err := actor.handleCommand(func() ([]*cloudevents.Event, error) {
		// Validate command
		if err := actor.validate.Struct(cmd); err != nil {
			return nil, err
		}
		return []*cloudevents.Event{
			event.CloudEvent("client", event.TypeClientRegisteredV1, &event.ClientRegistered{Id: actor.ID(), Salutation: cmd.Salutation, Name: cmd.Name, Surname: cmd.Surname, Phone: cmd.Phone, Email: cmd.Email}),
		}, nil
	})
	if err != nil {
		log.Error().Str("id", actor.ID()).Err(err).Msg("executing command")
	}
	return res, err
}

// applies the events to update the actor state
func (actor *Client) Apply(ces ...*cloudevents.Event) error {
	for _, ce := range ces {
		switch ce.Type() {
		case event.TypeClientRegisteredV1:
			var cr event.ClientRegistered
			if err := ce.DataAs(&cr); err != nil {
				return err
			}
			actor.state.Salutation = cr.Salutation
			actor.state.Surname = cr.Surname
			actor.state.Name = cr.Name
			actor.state.Phone = cr.Phone
			actor.state.Email = cr.Email
			actor.state.Version++
		default:
			return fmt.Errorf("unknown event type '%s'", ce.Type())
		}
	}
	return nil
}

// create new instances of an client actor
func ClientActorFactory() actor.Server {
	return &Client{
		validate: validator.New(),
	}
}
