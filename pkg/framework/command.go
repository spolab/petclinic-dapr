package framework

import (
	"github.com/rs/zerolog/log"
)

// A command handler is a function capable of a command against a state of type S
type CommandHandler func() error

type CommandExecutionLifecycle interface {
	Execute(actor EventSourcedAggregate, handle CommandHandler) error
}

type EventSourcedCommandExecutor struct {
	repo EventSourcedRepository
}

func (e EventSourcedCommandExecutor) Execute(aggregate EventSourcedAggregate, handle CommandHandler) error {
	log.Info().Str("id", aggregate.ID()).Msg("begin handleCommand")
	log.Debug().Str("id", aggregate.ID()).Msg("load event stream")
	err := e.repo.Load(aggregate)
	if err != nil {
		return err
	}
	//
	// Checks the invariants before executing the command
	//
	if err := aggregate.Check(); err != nil {
		return err
	}
	//
	// Execute the command
	//
	log.Debug().Str("id", aggregate.ID()).Msg("execute command")
	if err := handle(); err != nil {
		return err
	}
	//
	// Update the status using the events returned by the handler
	//
	log.Debug().Str("id", aggregate.ID()).Msg("apply events to state")
	if err := aggregate.Apply(aggregate.UncommittedEvents()...); err != nil {
		return err
	}
	//
	// Check invariants again
	//
	log.Debug().Str("id", aggregate.ID()).Msg("check invariants (post-execution")
	if err := aggregate.Check(); err != nil {
		return err
	}
	//
	// Append and store the events to the stream
	//
	log.Debug().Str("id", aggregate.ID()).Msg("store event stream")
	if err := e.repo.Save(aggregate); err != nil {
		return err
	}
	log.Info().Str("id", aggregate.ID()).Msg("end handleCommand")
	return nil
}
