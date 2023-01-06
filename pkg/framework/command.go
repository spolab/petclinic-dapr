/*
Copyright 2022 Alessandro Santini

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package framework

import (
	"github.com/rs/zerolog/log"
)

// A command handler is a function capable of a command against a state of type S
type CommandHandler func() error

type CommandExecutionLifecycle interface {
	Execute(actor EventSourcedAggregate, handle CommandHandler) error
}

type EventSourcedCommandLifecycle struct {
	Repository EventSourcedRepository
}

func (e EventSourcedCommandLifecycle) Execute(aggregate EventSourcedAggregate, handle CommandHandler) error {
	log.Info().Str("id", aggregate.ID()).Msg("begin handleCommand")
	//
	// It is worth reminding why the clear of the uncommitted events happens here:
	// Dapr's ActorManager saves the state of the actor only at the exit of the message handler.
	// We have no visibility, therefore, about whether those events have been truly committed or not.
	// What really matters, instead, is that we reload the actor fresh prior to command execution.
	// If the StateManager fails to persist the events, then an error will be returned to the caller,
	// so there will be no visibility of the events. The API therefore won't broadcast them.
	//
	log.Debug().Str("id", aggregate.ID()).Msg("clearing committed events")
	aggregate.ClearEvents()
	log.Debug().Str("id", aggregate.ID()).Msg("load event stream")
	if err := e.Repository.Load(aggregate); err != nil {
		return err
	}
	//
	// Checks the invariants before executing the command
	//
	log.Debug().Str("id", aggregate.ID()).Msg("check invariants (pre-execution)")
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
	log.Debug().Str("id", aggregate.ID()).Msg("check invariants (post-execution)")
	if err := aggregate.Check(); err != nil {
		return err
	}
	//
	// Append and store the events to the stream
	//
	log.Debug().Str("id", aggregate.ID()).Msg("store event stream")
	if err := e.Repository.Save(aggregate); err != nil {
		return err
	}
	log.Info().Str("id", aggregate.ID()).Msg("end handleCommand")
	return nil
}
