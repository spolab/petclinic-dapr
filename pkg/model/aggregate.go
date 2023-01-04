package model

import (
	"fmt"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/dapr/go-sdk/actor"
	"github.com/rs/zerolog/log"
	"github.com/spolab/petstore/pkg/command"
)

type CommandHandler func() ([]*cloudevents.Event, error)

type AggregateRoot[T any] interface {
	ID() string
	// Applies an event to alter the state.
	Apply(...*cloudevents.Event) error
	// Returns a pointer to the state. Please be aware this is mutable!
	State() T
	// Reconstructs the state from the event stream. Returns false if the aggregate does not exist.
	Load() (bool, error)
	// Tests the invariants and returns an error if the aggregate is inconsistent
	CheckInvariants() error
}

type BaseAggregateRoot[T any] struct {
	actor.ServerImplBase                      // DAPR actor framework
	state                T                    // cached copy of the state
	events               []*cloudevents.Event // cache of all the events loaded
	exists               bool                 // true if already cached in memory
}

func (base *BaseAggregateRoot[T]) State() T {
	return base.state
}

// Applies the events to the aggregate state
func (base *BaseAggregateRoot[T]) Apply(...*cloudevents.Event) error {
	return fmt.Errorf("not implemented")
}

// Invariants will be checked before and after the execution of a command
func (base *BaseAggregateRoot[T]) CheckInvariants() error {
	return nil
}

// handleCommand is a utility method that does covers the lifecycle of a command execution.
func (base *BaseAggregateRoot[T]) handleCommand(handler CommandHandler) (*command.ActorResponse, error) {
	log.Info().Str("id", base.ID()).Msg("begin handleCommand")
	log.Debug().Str("id", base.ID()).Msg("load event stream")
	if err := base.Load(); err != nil {
		return nil, err
	}
	//
	// Checks invariants before execution of the command (only if the aggregate exists)
	//
	if base.exists {
		log.Debug().Str("id", base.ID()).Msg("check invariants (pre-execution)")
		if err := base.CheckInvariants(); err != nil {
			return nil, err
		}
	}
	//
	// Execute the command
	//
	log.Debug().Str("id", base.ID()).Msg("execute command")
	events, err := handler()
	if err != nil {
		return nil, err
	}
	//
	// Update the status using the events returned by the handler
	//
	log.Debug().Str("id", base.ID()).Msg("apply events to state")
	if err := base.Apply(events...); err != nil {
		return nil, err
	}
	//
	// Check invariants again
	//
	log.Debug().Str("id", base.ID()).Msg("check invariants (post-execution")
	if err := base.CheckInvariants(); err != nil {
		return nil, err
	}
	//
	// Append and store the events to the stream
	//
	log.Debug().Str("id", base.ID()).Msg("store event stream")
	base.events = append(base.events, events...)
	if err := base.GetStateManager().Set(KeyEvents, base.events); err != nil {
		return nil, err
	}
	log.Info().Str("id", base.ID()).Msg("end handleCommand")
	return &command.ActorResponse{Status: command.StatusOK, Message: "", Events: events}, nil
}

// Load the state of the aggregate from the events log. Returns
func (base *BaseAggregateRoot[T]) Load() error {
	if base.exists {
		return nil
	}
	found, err := base.GetStateManager().Contains(KeyEvents)
	if err != nil {
		return err
	}
	if !found {
		return nil
	}
	err = base.GetStateManager().Get(KeyEvents, &base.events)
	if err != nil {
		return err
	}
	for _, event := range base.events {
		if err := base.Apply(event); err != nil {
			//
			// Erase the state if an error occurs.
			// In this way, any further attempt at using this actor instance will fail until the event stream handling gets fixed.
			//
			return err
		}
	}
	base.exists = true
	return nil
}
