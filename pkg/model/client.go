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
	"github.com/spolab/petstore/pkg/framework"
)

type Client struct {
	framework.BaseEventSourcedAggregate
	validate   *validator.Validate
	Salutation string
	Surname    string
	Name       string
	Phone      string
	Email      string
}

func (actor *Client) Type() string {
	return "client"
}

// register a new client
func (actor *Client) Register(ctx context.Context, cmd *command.RegisterClientCommand) (*command.ActorResponse, error) {
	err := actor.Lifecycle.Execute(actor, func() error {
		// The actor already exists
		if actor.Version == 0 {
			return fmt.Errorf("client id '%s' already exists", actor.ID())
		}
		// Validate command
		if err := actor.validate.Struct(cmd); err != nil {
			return err
		}
		actor.AppendEvent(event.CloudEvent("client", event.TypeClientRegisteredV1, &event.ClientRegistered{Id: actor.ID(), Salutation: cmd.Salutation, Name: cmd.Name, Surname: cmd.Surname, Phone: cmd.Phone, Email: cmd.Email}))
		return nil
	})
	if err != nil {
		log.Error().Str("id", actor.ID()).Err(err).Msg("executing command")
	}
	return nil, err
}

func (actor *Client) Apply(ces ...*cloudevents.Event) error {
	for _, ce := range ces {
		switch ce.Type() {
		case event.TypeClientRegisteredV1:
			var ev event.ClientRegistered
			if err := ce.DataAs(&ev); err != nil {
				return err
			}
			actor.Email = ev.Email
			actor.Name = ev.Name
			actor.Phone = ev.Phone
			actor.Salutation = ev.Salutation
			actor.Surname = ev.Surname
			actor.Version = 1
		}
	}
	return nil
}

func (actor *Client) Check() error {
	return nil
}

// create new instances of an client actor
func ClientActorFactory() actor.Server {
	result := &Client{
		validate: validator.New(),
	}
	return result
}
