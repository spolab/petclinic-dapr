package model

import (
	"context"
	"fmt"
	"time"

	"github.com/dapr/go-sdk/actor"
	"github.com/rs/zerolog/log"
	"github.com/spolab/petstore/pkg/command"
	"github.com/spolab/petstore/pkg/event"
)

// A data object containing the client contact details
type Client struct {
	Id           string                  `json:"id"`
	Salutation   string                  `json:"salutation"`
	Surname      string                  `json:"surname"`
	Name         string                  `json:"name"`
	Phone        string                  `json:"phone"`
	Email        string                  `json:"email"`
	Pets         map[string]*Pet         `json:"pets"`
	Appointments map[string]*Appointment `json:"appointments"`
}

type Pet struct {
	Id    string
	Name  string
	Breed string
}

type Appointment struct {
	Id    string    `json:"id"`
	Date  time.Time `json:"time"`
	PetId string    `json:"pet_id"`
}

type ClientActor struct {
	AggregateActor[Client]
}

func (actor *ClientActor) Type() string {
	return "client"
}

// register a new client
func (actor *ClientActor) Register(ctx context.Context, cmd *command.RegisterClientCommand) (Events, error) {
	return actor.commandLifecycle(actor, cmd, func() error {
		if actor.exists {
			return fmt.Errorf("client id '%s' already exists", actor.ID())
		}
		// Append the events
		cr := &event.ClientRegistered{
			Id:         actor.ID(),
			Salutation: cmd.Salutation,
			Name:       cmd.Name,
			Surname:    cmd.Surname,
			Phone:      cmd.Phone,
			Email:      cmd.Email,
		}
		// Generate the stream of events and apply them
		events := Events{
			event.CloudEvent(event.FromSource("client"), event.OfType(event.TypeClientRegisteredV1), event.WithDataAsJSON(cr)),
		}
		actor.Append(events...)
		return nil
	})
}

func (actor *ClientActor) RegisterPet(ctx context.Context, cmd *command.RegisterPetCommand) (Events, error) {
	return actor.commandLifecycle(actor, cmd, func() error {
		//
		// Return an error if the client does not exist
		//
		if !actor.exists {
			return fmt.Errorf("client does not exist")
		}
		//
		// Return an error if a pet with the same chip ID exists
		//
		_, exists := actor.snapshot.Pets[cmd.ChipID]
		if exists {
			return fmt.Errorf("a pet with such ID already exists")
		}
		//
		// All good, append the event
		//
		ce := event.CloudEvent(event.FromSource("client"), event.OfType(event.TypePetRegisteredV1), event.WithDataAsJSON(&event.PetRegistered{
			ClientID: actor.ID(),
			ChipID:   cmd.ChipID,
			Name:     cmd.Name,
			Breed:    cmd.Breed,
		}))
		actor.Append(ce)
		return nil
	})
}

func (actor *ClientActor) RequestVisit(ctx context.Context, cmd *command.RequestAppointmentCommand) (Events, error) {
	return actor.commandLifecycle(actor, cmd, func() error {
		if !actor.exists {
			return fmt.Errorf("no such pet exists")
		}
		ar := event.CloudEvent(event.FromSource("client"), event.OfType(event.TypeAppointmentRequestedV1), event.WithDataAsJSON(event.AppointmentRequested{
			ClientID:  actor.ID(),
			ChipID:    cmd.ChipID,
			Specialty: cmd.Specialty,
		}))
		actor.Append(ar)
		return nil
	})
}

// Applies the events without updating the state of the aggregate
func clientMutator(state *Client, events Events) error {
	log.Info().Str("id", state.Id).Msg("begin apply")
	for _, ce := range events {
		switch ce.Type() {
		case event.TypeClientRegisteredV1:
			var data event.ClientRegistered
			if err := ce.DataAs(&data); err != nil {
				return err
			}
			state.Id = data.Id
			state.Salutation = data.Salutation
			state.Surname = data.Surname
			state.Name = data.Name
			state.Phone = data.Phone
			state.Email = data.Email
			state.Pets = make(map[string]*Pet)
		case event.TypePetRegisteredV1:
			var pr event.PetRegistered
			if err := ce.DataAs(&pr); err != nil {
				return err
			}
			pet := &Pet{
				Id:   pr.ChipID,
				Name: pr.Name,
			}
			state.Pets[pr.ChipID] = pet
		default:
			log.Warn().Str("id", state.Id).Str("type", ce.Type()).Msg("unknown event type")
		}
	}
	log.Info().Str("id", state.Id).Msg("end apply")
	return nil
}

// create new instances of an client actor
func ClientActorFactory() actor.Server {
	log.Info().Msg("begin clientactorfactory")
	result := &ClientActor{}
	result.snapshot = Client{}
	result.uncommittedEvents = Events{}
	result.mutator = clientMutator
	result.commandLifecycle = DefaultCommandLifecycle[Client]
	log.Info().Msg("end clientactorfactory")
	return result
}
