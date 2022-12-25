package server

import (
	"context"
	"fmt"

	"github.com/spolab/petstore/pkg/command"
	"github.com/spolab/petstore/pkg/event"
	"go.uber.org/zap"
)

const keyContactState = "contact"

func (a *OwnerActor) Register(ctx context.Context, cmd *command.RegisterOwner) error {
	a.logger.Info("start registering new owner", zap.String("id", a.ID()))
	// validate the command
	if err := a.validate.Struct(cmd); err != nil {
		// This should return a 400 instead of a 500, how do we do that?
		a.logger.Error("validating command", zap.String("id", a.ID()), zap.Error(err))
		return err
	}
	// exit if the owner already exists
	found, err := a.GetStateManager().Contains(keyContactState)
	if err != nil {
		a.logger.Error("checking if the owner already exists", zap.String("id", a.ID()), zap.Error(err))
		return err
	}
	if found {
		a.logger.Debug("owner already exists", zap.String("id", a.ID()))
		return fmt.Errorf("owner '%s' is already registered", a.ID())
	}
	// store the state of the owner
	state := OwnerState{
		Salutation: cmd.Salutation,
		Name:       cmd.Name,
		Surname:    cmd.Surname,
	}
	if err := a.GetStateManager().Set(keyContactState, &state); err != nil {
		a.logger.Error("setting state", zap.String("id", a.ID()), zap.Error(err))
		return err
	}
	// all good, publish the event
	if err := a.dapr.PublishEvent(ctx, a.broker, a.topic, event.OwnerRegistered{
		ID:         a.ID(),
		Salutation: cmd.Salutation,
		Name:       cmd.Name,
		Surname:    cmd.Surname,
	}); err != nil {
		a.logger.Error("publishing event", zap.String("id", a.ID()), zap.Error(err))
		return err
	}
	return nil
}
