package server

import (
	"context"
	"fmt"

	"github.com/spolab/petclinic/owner/command"
	"go.uber.org/zap"
)

const keyStateContact = "contact"

func (a *OwnerActor) Register(ctx context.Context, cmd *command.RegisterOwner) (*OwnerState, error) {
	a.logger.Info("start registering new owner", zap.String("id", a.ID()))
	found, err := a.GetStateManager().Contains(keyStateContact)
	if found {
		return nil, fmt.Errorf("owner %s is already registered", a.ID())
	}
	if err != nil {
		return nil, err
	}
	state := OwnerState{
		ID:         a.ID(),
		Salutation: cmd.Salutation,
		Name:       cmd.Name,
		Surname:    cmd.Surname,
	}
	if err := a.GetStateManager().Set(keyStateContact, state); err != nil {
		return nil, err
	}
	return &state, nil
}
