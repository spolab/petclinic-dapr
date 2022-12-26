package owner

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spolab/petstore/gen/api"
	"github.com/spolab/petstore/pkg/event"
)

const keyContactState = "contact"

func (a *OwnerActor) Register(ctx context.Context, cmd *api.RegisterOwnerCommand) error {
	log.Info().Str("id", a.ID()).Msg("start registering new owner")
	// validate the command
	if err := a.validate.Struct(cmd); err != nil {
		// This should return a 400 instead of a 500, how do we do that?
		log.Error().Err(err).Str("id", a.ID()).Msg("validating command")
		return err
	}
	// exit if the owner already exists
	found, err := a.GetStateManager().Contains(keyContactState)
	if err != nil {
		log.Error().Str("id", a.ID()).Err(err).Msg("checking if the owner already exists")
		return err
	}
	if found {
		log.Debug().Str("id", a.ID()).Msg("owner already exists")
		return fmt.Errorf("owner '%s' is already registered", a.ID())
	}
	// store the state of the owner
	state := OwnerState{
		Salutation: cmd.Salutation.String(),
		Name:       cmd.Name,
		Surname:    cmd.Surname,
	}
	if err := a.GetStateManager().Set(keyContactState, &state); err != nil {
		log.Error().Str("id", a.ID()).Err(err).Msg("setting state")
		return err
	}
	// all good, publish the event
	if err := a.dapr.PublishEvent(ctx, a.broker, a.topic, event.OwnerRegistered{
		ID:         a.ID(),
		Salutation: cmd.Salutation.String(),
		Name:       cmd.Name,
		Surname:    cmd.Surname,
	}); err != nil {
		log.Error().Str("id", a.ID()).Err(err).Msg("publishing event")
		return err
	}
	return nil
}
