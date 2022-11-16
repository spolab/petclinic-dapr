package model

import (
	"context"
	"fmt"
	"reflect"

	"github.com/spolab/petclinic/owner/pkg/api"
)

type Owner struct {
	Id      string
	Version string
	Events  []any
	State   struct {
		Salutation string
		Surname    string
		Name       string
	}
}

// Apply is a method that alters the state of the aggregate based on the event provided
func (o *Owner) Apply(ctx context.Context, rawEvent any) error {
	switch event := rawEvent.(type) {
	case api.OwnerRegistered:
		o.Id = event.Id
		o.State.Salutation = event.Salutation
		o.State.Surname = event.Surname
		o.State.Name = event.Name
	default:
		return fmt.Errorf("unknown event type '%s'", reflect.TypeOf(rawEvent).Name())
	}
	o.Events = append(o.Events, rawEvent)
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
	if err := s.Apply(ctx, api.OwnerRegistered{
		Owner: cmd.Owner,
	}); err != nil {
		return err
	}
	return nil
}
