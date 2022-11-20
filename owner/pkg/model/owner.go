package model

import (
	"context"
	"fmt"
	"reflect"
	"time"

	ce "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	"github.com/spolab/petclinic/owner/pkg/api"
)

const OwnerRegisteredType = "spolab/petclinic/OwnerRegistered/v1"

type Owner struct {
	Id      string `json:"id"`
	Version string `json:"version"`
	// TODO better to use a pointer to events or a copy of the events?
	UncommittedEvents []ce.Event
	State             struct {
		Salutation string `json:"salutation"`
		Surname    string `json:"surname"`
		Name       string `json:"name"`
		Phone      string `json:"phone"`
		Email      string `json:"email"`
	} `json:"state"`
}

// Apply is a method that alters the state of the aggregate based on the event provided
// TODO do we need context here? there should never be any input/output here
func (o *Owner) Apply(ctx context.Context, event ce.Event) error {
	switch event.Type() {
	//
	case OwnerRegisteredType:
		var data api.OwnerRegistered
		if err := event.DataAs(&data); err != nil {
			return err
		}
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
func (o *Owner) Append(ctx context.Context, event ce.Event) error {
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
	return s.Append(ctx, BuildEvent(OwnerRegisteredType, "owner", event))
}

func BuildEvent(kind string, source string, data any) ce.Event {
	result := ce.NewEvent()
	result.SetSpecVersion(ce.VersionV1)
	result.SetID(uuid.NewString())
	result.SetType(kind)
	result.SetTime(time.Now())
	result.SetSource(source)
	result.SetDataContentEncoding(ce.EncodingStructured.String())
	result.SetData(ce.ApplicationJSON, data)
	return result
}
