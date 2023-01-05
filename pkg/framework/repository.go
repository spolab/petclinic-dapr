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
		if err := target.GetStateManager().Get(daprKeyEvents, events); err != nil {
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
