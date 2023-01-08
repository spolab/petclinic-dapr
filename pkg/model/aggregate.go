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
)

const (
	keyUncommittedEvents = "uncommittedEvent"
)

type AggregateActor interface {
	// Mutates the state of the aggregate by applying the given events
	Apply(...*cloudevents.Event) error
}

// An event sourced actor is capable of dealing with an event source state
type BaseAggregateActor struct {
	actor.ServerImplBase
	uncommittedEvents []*cloudevents.Event
}

func (a BaseAggregateActor) Append(events ...*cloudevents.Event) {
	a.uncommittedEvents = append(a.uncommittedEvents, events...)
	a.GetStateManager().Set(keyClientEvents, a.uncommittedEvents)
}

// Returns the queue of events that have not been stored/propagated yet
func (a BaseAggregateActor) UncommittedEvents() ([]*cloudevents.Event, error) {
	return a.uncommittedEvents, nil
}

// Clears the queue of uncommitted events
func (a BaseAggregateActor) ClearEvents() error {
	a.uncommittedEvents = []*cloudevents.Event{}
	return a.GetStateManager().Set(keyUncommittedEvents, a.uncommittedEvents)
}
