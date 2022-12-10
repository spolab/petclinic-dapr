package server

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dapr/go-sdk/service/common"
	"github.com/spolab/petclinic/owner/command"
	"github.com/spolab/petclinic/owner/event"
	"go.uber.org/zap"
)

const keyContactState = "contact"

func (a *OwnerActor) Register(ctx context.Context, in *common.InvocationEvent) (*common.Content, error) {
	a.logger.Info("start registering new owner", zap.String("id", a.ID()))
	// unmarshal the request and turn it into a command
	var cmd command.RegisterOwner
	if err := json.Unmarshal(in.Data, &cmd); err != nil {
		return nil, err
	}
	// validate the command
	if err := a.validate.Struct(cmd); err != nil {
		a.logger.Error("validating command", zap.String("id", a.ID()), zap.Error(err))
		return nil, err
	}
	// exit if the owner already exists
	found, err := a.GetStateManager().Contains(keyContactState)
	if err != nil {
		a.logger.Error("checking if the owner already exists", zap.String("id", a.ID()), zap.Error(err))
		return nil, err
	}
	if found {
		a.logger.Debug("owner already exists", zap.String("id", a.ID()))
		return nil, fmt.Errorf("owner '%s' is already registered", a.ID())
	}
	// store the state of the owner
	state := OwnerState{
		Salutation: cmd.Salutation,
		Name:       cmd.Name,
		Surname:    cmd.Surname,
	}
	if err := a.GetStateManager().Set(keyContactState, &state); err != nil {
		a.logger.Error("setting state", zap.String("id", a.ID()), zap.Error(err))
		return nil, err
	}
	// all good, publish the event
	if err := a.dapr.PublishEvent(ctx, a.broker, a.topic, event.OwnerRegistered{
		ID:         a.ID(),
		Salutation: cmd.Salutation,
		Name:       cmd.Name,
		Surname:    cmd.Surname,
	}); err != nil {
		a.logger.Error("publishing event", zap.String("id", a.ID()), zap.Error(err))
		return nil, err
	}
	// return the response
	a.logger.Info("end registering new owner", zap.String("id", a.ID()))
	return JSON(state)
}
