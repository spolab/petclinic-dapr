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
)

type EventSourcedRepository interface {
	Load(EventSourcedAggregate) error
	// Saves the aggregate's uncommitted events
	Save(EventSourcedAggregate) error
}

const (
	daprKeyEvents = "events"
)

type EventSourcedActorRepository struct {
}

func (repo EventSourcedActorRepository) Load(target EventSourcedAggregate) error {
	events := []*cloudevents.Event{}
	found, err := target.GetStateManager().Contains(daprKeyEvents)
	if err != nil {
		return err
	}
	if found {
		if err := target.GetStateManager().Get(daprKeyEvents, &events); err != nil {
			return err
		}
		if err := target.Apply(events...); err != nil {
			return err
		}
	}
	return nil
}

func (repo EventSourcedActorRepository) Save(source EventSourcedAggregate) error {
	if err := source.GetStateManager().Set(daprKeyEvents, source.UncommittedEvents()); err != nil {
		return err
	}
	source.ClearEvents()
	return nil
}
