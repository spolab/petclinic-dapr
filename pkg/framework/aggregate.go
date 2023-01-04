package framework

import (
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/dapr/go-sdk/actor"
	"github.com/rs/zerolog/log"
	"github.com/spolab/petstore/pkg/command"
)

const (
	daprKeyEvents = "events"
)

// A command handler is a function capable of a command against a state of type S
type CommandHandler func() ([]*cloudevents.Event, error)

type EventSourcedState interface {
	// EventApplicator is responsible for applying a stream of events to a state of type S
	Apply(events ...*cloudevents.Event) error
	// Check the aggregate invariantes
	Check() error
}

// An event sourced actor is capable of dealing with an event source state
type EventSourcedActor[T EventSourcedState] struct {
	actor.ServerImplBase
	State  T
	Stream []*cloudevents.Event
}

func (actor *EventSourcedActor[T]) Load() error {
	var events []*cloudevents.Event
	var state T
	found, err := actor.GetStateManager().Contains(daprKeyEvents)
	if err != nil {
		return err
	}
	if found {
		if err := actor.GetStateManager().Get(daprKeyEvents, &events); err != nil {
			return err
		}
		if err := state.Apply(events...); err != nil {
			return err
		}
	}
	actor.State = state
	return nil
}

// Stores the Dapr actor events using the state manager.
// The ID is ignored because the state manager already knows it (it is stateful).
func (actor EventSourcedActor[T]) Save() error {
	if err := actor.GetStateManager().Set(daprKeyEvents, actor.Stream); err != nil {
		return err
	}
	return nil
}

// handleCommand is a utility method that does covers the lifecycle of a command execution.
func (base EventSourcedActor[T]) HandleCommand(cmd any, handle CommandHandler) (*command.ActorResponse, error) {
	log.Info().Str("id", base.ID()).Msg("begin handleCommand")
	log.Debug().Str("id", base.ID()).Msg("load event stream")
	err := base.Load()
	if err != nil {
		return nil, err
	}
	//
	// Checks the invariants before executing the command
	//
	if err := base.State.Check(); err != nil {
		return nil, err
	}
	//
	// Execute the command
	//
	log.Debug().Str("id", base.ID()).Msg("execute command")
	events, err := handle()
	if err != nil {
		return nil, err
	}
	//
	// Update the status using the events returned by the handler
	//
	log.Debug().Str("id", base.ID()).Msg("apply events to state")
	if err := base.State.Apply(events...); err != nil {
		return nil, err
	}
	//
	// Check invariants again
	//
	log.Debug().Str("id", base.ID()).Msg("check invariants (post-execution")
	if err := base.State.Check(); err != nil {
		return nil, err
	}
	//
	// Append and store the events to the stream
	//
	log.Debug().Str("id", base.ID()).Msg("store event stream")
	base.Stream = append(base.Stream, events...)
	if err := base.Save(); err != nil {
		return nil, err
	}
	log.Info().Str("id", base.ID()).Msg("end handleCommand")
	return &command.ActorResponse{Status: command.StatusOK, Message: "", Events: events}, nil
}
