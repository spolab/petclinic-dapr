package model

import (
	"context"
	"fmt"
	"reflect"

	"github.com/spolab/petclinic/owner/pkg/api"
)

const OwnerRegisteredUrn = "spolab/petclinic/OwnerRegistered/v1"

type Owner struct {
	Id                string
	Version           string
	UncommittedEvents []any
	State             struct {
		Salutation string
		Surname    string
		Name       string
		Phone      string
		Email      string
	}
}

// Apply is a method that alters the state of the aggregate based on the event provided
func (o *Owner) Apply(ctx context.Context, event any) error {
	switch data := event.(type) {
	//
	case api.OwnerRegistered:
		o.Id = data.Id
		o.State.Salutation = data.Salutation
		o.State.Surname = data.Surname
		o.State.Name = data.Name
		o.State.Phone = data.Phone
		o.State.Email = data.Email
	//
	default:
		return fmt.Errorf("unknown event type '%s'", reflect.TypeOf(event))
	}
	return nil
}

// Append applies a new event to the aggregate root and appends it to the list of uncommitted events
func (o *Owner) Append(ctx context.Context, event any) error {
	err := o.Apply(ctx, event)
	if err == nil {
		o.UncommittedEvents = append(o.UncommittedEvents, event)
	}
	return err
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
	event := api.OwnerRegistered{
		Id:         cmd.Id,
		Salutation: cmd.Salutation,
		Surname:    cmd.Surname,
		Name:       cmd.Name,
		Phone:      cmd.Phone,
		Email:      cmd.Email,
	}
	return s.Append(ctx, event)
}
