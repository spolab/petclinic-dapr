package server

import (
	"context"
	"fmt"

	"github.com/spolab/petclinic/owner/command"
	"github.com/spolab/petclinic/owner/event"
	"go.uber.org/zap"
)

const keyStateContact = "contact"

func (a *OwnerActor) Register(ctx context.Context, cmd *command.RegisterOwner) (*OwnerState, error) {
	a.logger.Info("start registering new owner", zap.String("id", a.ID()))
	found, err := a.GetStateManager().Contains(keyStateContact)
	if err != nil {
		a.logger.Error("checking if the owner already exists", zap.String("id", a.ID()), zap.Error(err))
		return nil, err
	}
	if found {
		a.logger.Debug("owner already exists", zap.String("id", a.ID()))
		return nil, fmt.Errorf("owner '%s' is already registered", a.ID())
	}
	if err := a.validate.Struct(cmd); err != nil {
		a.logger.Error("validating command", zap.String("id", a.ID()), zap.Error(err))
		return nil, err
	}
	state := OwnerState{
		Salutation: cmd.Salutation,
		Name:       cmd.Name,
		Surname:    cmd.Surname,
	}
	if err := a.GetStateManager().Set(keyStateContact, &state); err != nil {
		a.logger.Error("setting state", zap.String("id", a.ID()), zap.Error(err))
		return nil, err
	}
	if err := a.dapr.PublishEvent(ctx, a.broker, a.topic, event.OwnerRegistered{
		ID:         a.ID(),
		Salutation: cmd.Salutation,
		Name:       cmd.Name,
		Surname:    cmd.Surname,
	}); err != nil {
		a.logger.Error("publishing event", zap.String("id", a.ID()), zap.Error(err))
		return nil, err
	}
	a.logger.Info("end registering new owner", zap.String("id", a.ID()))
	return &state, nil
}
