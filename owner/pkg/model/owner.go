package model

import (
	"context"
	"fmt"

	ce "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	"github.com/spolab/petclinic/owner/pkg/api"
)

const OwnerRegisteredUrn = "spolab/petclinic/OwnerRegistered/v1"

type Owner struct {
	Id                string
	Version           string
	UncommittedEvents []ce.Event
	State             struct {
		Salutation string
		Surname    string
		Name       string
	}
}

// Apply is a method that alters the state of the aggregate based on the event provided
func (o *Owner) Apply(ctx context.Context, event ce.Event) error {
	switch event.DataSchema() {
	//
	case OwnerRegisteredUrn:
		var data api.OwnerRegistered
		event.DataAs(&data)
		o.Id = data.Id
		o.State.Salutation = data.Salutation
		o.State.Surname = data.Surname
		o.State.Name = data.Name
	//
	default:
		return fmt.Errorf("unknown event type '%s'", event.DataSchema())
	}
	return nil
}

// Register a new owner.
func (s *Owner) Register(ctx context.Context, cmd *api.RegisterOwner) error {
	if cmd == nil {
		return fmt.Errorf("command cannot be nil")
	}
	if err := cmd.Validate(); err != nil {
		return err
	}
	// Creates the cloudevent and applies it
	// Apply the event and add it to the queue of uncommitted events
	event := NewOwnerCreated(cmd.Id, cmd.Salutation, cmd.Surname, cmd.Name)
	if err := s.Apply(ctx, event); err != nil {
		return err
	}
	s.UncommittedEvents = append(s.UncommittedEvents, event)
	return nil
}

func NewOwnerCreated(id string, salutation string, surname string, name string) ce.Event {
	result := ce.NewEvent()
	result.SetID(uuid.New().String())
	result.SetDataSchema(OwnerRegisteredUrn)
	result.SetData(ce.ApplicationJSON, api.OwnerRegistered{
		Id:         id,
		Salutation: salutation,
		Surname:    surname,
		Name:       name,
	})
	return result
}
