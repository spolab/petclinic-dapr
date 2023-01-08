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
	keyClientContact = "contact"
	keyClientPets    = "pets"
	keyClientEvents  = "events"
)

// A data object containing the client contact details
type ClientState struct {
	Salutation string `json:"salutation"`
	Surname    string `json:"surname"`
	Name       string `json:"name"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
}

type PetState struct {
	Name   string
	ChipID string
	Breed  string
}

type ClientActor struct {
	BaseAggregateActor
	validate *validator.Validate
	contact  *ClientState
	pets     map[string]*PetState
	version  int
}

func (actor *ClientActor) Type() string {
	return "client"
}

// register a new client
func (actor *ClientActor) Register(ctx context.Context, cmd *command.RegisterClientCommand) ([]*cloudevents.Event, error) {
	exists, err := actor.GetStateManager().Contains(keyClientContact)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("client id '%s' already exists", actor.ID())
	}
	// Validate command
	if err := actor.validate.Struct(cmd); err != nil {
		return nil, err
	}
	// Append the events
	cr := &event.ClientRegistered{
		Id:         actor.ID(),
		Salutation: cmd.Salutation,
		Name:       cmd.Name,
		Surname:    cmd.Surname,
		Phone:      cmd.Phone,
		Email:      cmd.Email,
		Version:    actor.version + 1,
	}
	// Generate the stream of events and apply them
	events := []*cloudevents.Event{
		event.CloudEvent(event.FromSource("client"), event.OfType(event.TypeClientRegisteredV1), event.WithDataAsJSON(cr)),
	}
	if err := actor.Apply(events...); err != nil {
		log.Error().Str("id", actor.ID()).Err(err).Msg("executing command")
		return nil, err
	}
	actor.Append(events...)
	return actor.UncommittedEvents()
}

func (actor *ClientActor) RegisterPet(ctx context.Context, cmd *command.RegisterPetCommand) error {
	//
	// Validate the command
	//
	err := actor.validate.Struct(cmd)
	if err != nil {
		return err
	}
	//
	// Return an error if the client does not exist
	//
	exist, err := actor.GetStateManager().Contains(keyClientPets)
	if err != nil {
		return err
	}
	if !exist {
		return fmt.Errorf("client does not exist")
	}
	err = actor.GetStateManager().Get(keyClientPets, &actor.pets)
	if err != nil {
		return err
	}
	//
	// Return an error if a pet with the same chip ID exists
	//
	_, exists := actor.pets[cmd.ChipID]
	if exists {
		return fmt.Errorf("a pet with such ID already exists")
	}
	//
	// All good, emit the event and process it
	//
	ce := event.CloudEvent(event.FromSource("client"), event.OfType(event.TypePetRegisteredV1), event.WithDataAsJSON(&event.PetRegistered{
		ClientID: actor.ID(),
		ChipID:   cmd.ChipID,
		Name:     cmd.Name,
		Breed:    cmd.Breed,
	}))
	if err := actor.Apply(ce); err != nil {
		return err
	}
	actor.Append(ce)
	if err := actor.Save(); err != nil {
		return err
	}
	return nil
}

func (actor *ClientActor) RequestVisit(ctx context.Context, cmd *command.RequestAppointmentCommand) ([]*cloudevents.Event, error) {
	var events []*cloudevents.Event
	err := actor.GetStateManager().Get(keyClientPets, &actor.pets)
	if err != nil {
		return nil, err
	}

	return events, nil
}

// Applies the events without updating the state of the aggregate
func (actor *ClientActor) Apply(ces ...*cloudevents.Event) error {
	log.Info().Str("id", actor.ID()).Msg("begin apply")
	for _, ce := range ces {
		switch ce.Type() {
		case event.TypeClientRegisteredV1:
			var cr event.ClientRegistered
			if err := ce.DataAs(&cr); err != nil {
				return err
			}
			actor.contact = &ClientState{
				Salutation: cr.Salutation,
				Surname:    cr.Surname,
				Name:       cr.Name,
				Phone:      cr.Phone,
				Email:      cr.Email,
			}
			actor.pets = make(map[string]*PetState)
		case event.TypePetRegisteredV1:
			var pr event.PetRegistered
			if err := ce.DataAs(&pr); err != nil {
				return err
			}
			pet := &PetState{
				ChipID: pr.ChipID,
				Name:   pr.Name,
			}
			actor.pets[pr.ChipID] = pet
		default:
			log.Warn().Str("id", actor.ID()).Str("type", ce.Type()).Msg("unknown event type")
		}
	}
	if err := actor.GetStateManager().Set("contact", actor.contact); err != nil {
		return err
	}
	log.Info().Str("id", actor.ID()).Msg("end apply")
	return nil
}

// Save is a utility method to persist the state of the client
func (actor *ClientActor) Save() error {
	if err := actor.GetStateManager().Set(keyClientContact, actor.contact); err != nil {
		return err
	}
	if err := actor.GetStateManager().Set(keyClientPets, actor.pets); err != nil {
		return err
	}
	return nil
}

// create new instances of an client actor
func ClientActorFactory() actor.Server {
	log.Info().Msg("begin clientactorfactory")
	result := &ClientActor{
		validate: validator.New(),
	}
	log.Info().Msg("end clientactorfactory")
	return result
}
