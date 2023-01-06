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
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/dapr/go-sdk/actor"
)

type EventSourcedAggregate interface {
	// Returns the ID of the aggregate
	ID() string
	// Mutates the state of the aggregate by applying the given events
	Apply(...*cloudevents.Event) error
	// Tells whether the state invariants hold, error otherwise
	Check() error
	// Returns the events that have not been committed since last save
	UncommittedEvents() []*cloudevents.Event
	// Adds events to the list of uncommitted events
	AppendEvent(...*cloudevents.Event)
	// Clears the queue of uncommitted events
	ClearEvents()
	// I really wish this was not here but it has to because of the way Dapr works
	GetStateManager() actor.StateManager
}

// An event sourced actor is capable of dealing with an event source state
type BaseEventSourcedAggregate struct {
	actor.ServerImplBase
	Lifecycle         CommandExecutionLifecycle
	uncommittedEvents []*cloudevents.Event
	Version           int
	Deleted           bool
}

// Sets a command execution lifecycle manager
func (a *BaseEventSourcedAggregate) SetLifecycle(lifecycle CommandExecutionLifecycle) {
	a.Lifecycle = lifecycle
}

// Returns the events that have not been committed since last save
func (a *BaseEventSourcedAggregate) UncommittedEvents() []*cloudevents.Event {
	return a.uncommittedEvents
}

// Adds events to the list of uncommitted events
func (a *BaseEventSourcedAggregate) AppendEvent(events ...*cloudevents.Event) {
	a.uncommittedEvents = append(a.uncommittedEvents, events...)
}

// Clears the queue of uncommitted events
func (a *BaseEventSourcedAggregate) ClearEvents() {
	a.uncommittedEvents = []*cloudevents.Event{}
}
