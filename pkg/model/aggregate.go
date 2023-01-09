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
package model

import (
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/dapr/go-sdk/actor"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
)

const (
	keySnapshot          = "snapshot"
	keyUncommittedEvents = "uncommittedEvents"
)

type Events []*cloudevents.Event

type Aggregate[T any] interface {
	// Returns the ID of the aggregate
	ID() string
	// Lists the events that have not been cleared yet
	UncommittedEvents() Events
	// Loads the aggregate from the persistence
	Load() error
	// Saves the aggregate
	Save() error
	// Mutates the state through the given events
	Apply(Events) error
}

type StateMutator[T any] func(*T, Events) error

type CommandLifecycle[T any] func(Aggregate[T], any, func() error) (Events, error)

// An event sourced actor is capable of dealing with an event source state
type AggregateActor[T any] struct {
	actor.ServerImplBase
	snapshot          T
	uncommittedEvents Events
	exists            bool
	mutator           StateMutator[T]
	commandLifecycle  CommandLifecycle[T]
}

func (a *AggregateActor[T]) Load() error {
	//
	// Load the snapshot
	//
	if exists, err := a.GetStateManager().Contains(keySnapshot); err != nil {
		return err
	} else {
		a.exists = exists
	}
	if a.exists {
		if err := a.GetStateManager().Get(keySnapshot, &a.snapshot); err != nil {
			return err
		}
	}
	//
	// Load the uncommitted events
	//
	if exists, err := a.GetStateManager().Contains(keyUncommittedEvents); err != nil {
		return err
	} else {
		a.exists = exists
	}
	if a.exists {
		if err := a.GetStateManager().Get(keyUncommittedEvents, &a.uncommittedEvents); err != nil {
			return err
		}
	}
	return nil
}

func (a *AggregateActor[T]) Save() error {
	if err := a.GetStateManager().Set(keySnapshot, a.snapshot); err != nil {
		return err
	}
	if err := a.GetStateManager().Set(keyUncommittedEvents, a.uncommittedEvents); err != nil {
		return err
	}
	return nil
}

func (a *AggregateActor[T]) Append(events ...*cloudevents.Event) {
	a.uncommittedEvents = append(a.uncommittedEvents, events...)
}

func (a *AggregateActor[T]) Apply(Events) error {
	if err := a.mutator(&a.snapshot, a.uncommittedEvents); err != nil {
		log.Error().Str("id", a.ID()).Err(err).Msg("mutating state")
		return err
	}
	return nil
}

// Returns the queue of events that have not been stored/propagated yet
func (a *AggregateActor[T]) UncommittedEvents() Events {
	return a.uncommittedEvents
}

// Clears the queue of uncommitted events
func (a *AggregateActor[T]) ClearEvents() error {
	a.uncommittedEvents = Events{}
	return a.GetStateManager().Set(keyUncommittedEvents, a.uncommittedEvents)
}

func DefaultCommandLifecycle[T any](a Aggregate[T], cmd any, handler func() error) (Events, error) {
	//
	// Validate the command
	// This can be optimized by caching the validator in the lifecycle
	//
	if err := validator.New().Struct(&cmd); err != nil {
		log.Error().Str("id", a.ID()).Err(err).Msg("validating command")
		return nil, err
	}
	//
	// Retrieve the actor state - assuming caching is implemented there
	//
	if err := a.Load(); err != nil {
		log.Error().Str("id", a.ID()).Err(err).Msg("retrieving state")
		return nil, err
	}
	//
	// Execute the command handler
	//
	if err := handler(); err != nil {
		log.Error().Str("id", a.ID()).Err(err).Msg("executing command")
		return nil, err
	}
	//
	// Apply the uncommitted events
	//
	if err := a.Apply(a.UncommittedEvents()); err != nil {
		log.Error().Str("id", a.ID()).Err(err).Msg("applying events")
		return nil, err
	}
	//
	// Save the aggregate
	//
	if err := a.Save(); err != nil {
		log.Error().Str("id", a.ID()).Err(err).Msg("saving state")
		return nil, err
	}
	return a.UncommittedEvents(), nil
}
