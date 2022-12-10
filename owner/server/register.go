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

const keyStateContact = "contact"

func (a *OwnerActor) Register(ctx context.Context, req *common.InvocationEvent) (*common.Content, error) {
	a.logger.Info("start registering new owner", zap.String("id", a.ID()))
	// deserialize command
	var cmd command.RegisterOwner
	err := json.Unmarshal(req.Data, &cmd)
	if err != nil {
		a.logger.Error("deserializing the request", zap.String("id", a.ID()), zap.Error(err))
		return nil, err
	}
	// command is parsed, validate it
	if err := a.validate.Struct(cmd); err != nil {
		a.logger.Error("validating command", zap.String("id", a.ID()), zap.Error(err))
		return nil, err
	}
	// command is valid, return an error if the owner already exists
	found, err := a.GetStateManager().Contains(keyStateContact)
	if err != nil {
		a.logger.Error("checking if the owner already exists", zap.String("id", a.ID()), zap.Error(err))
		return nil, err
	}
	if found {
		a.logger.Debug("owner already exists", zap.String("id", a.ID()))
		return nil, fmt.Errorf("owner '%s' is already registered", a.ID())
	}
	// owner does not exist, persist it
	state := OwnerState{
		Salutation: cmd.Salutation,
		Name:       cmd.Name,
		Surname:    cmd.Surname,
	}
	if err := a.GetStateManager().Set(keyStateContact, &state); err != nil {
		a.logger.Error("setting state", zap.String("id", a.ID()), zap.Error(err))
		return nil, err
	}
	// owner persisted correctly, emit the ownerregistered event
	if err := a.dapr.PublishEvent(ctx, a.broker, a.topic, event.OwnerRegistered{
		ID:         a.ID(),
		Salutation: cmd.Salutation,
		Name:       cmd.Name,
		Surname:    cmd.Surname,
	}); err != nil {
		a.logger.Error("publishing event", zap.String("id", a.ID()), zap.Error(err))
		return nil, err
	}
	// send a response back to the caller
	data, err := json.Marshal(&state)
	if err != nil {
		a.logger.Error("serializing response", zap.String("id", a.ID()), zap.Error(err))
		return nil, err
	}
	result := &common.Content{
		ContentType: "application/json",
		Data:        data,
	}
	a.logger.Info("end registering new owner", zap.String("id", a.ID()))
	return result, nil
}
