package framework

import (
	"testing"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/stretchr/testify/assert"
)

type TestAggregate struct {
	BaseEventSourcedAggregate
	t            *testing.T
	ApplyCount   int
	CheckCount   int
	ExecuteCount int
}

func (a *TestAggregate) Apply(...*cloudevents.Event) error {
	a.ApplyCount++
	return nil
}

func (a *TestAggregate) Check() error {
	a.CheckCount++
	return nil
}

func TestID(t *testing.T) {
	instance := TestAggregate{t: t}
	instance.SetID("id")
	assert.Equal(t, "id", instance.ID())
}

func TestEvents(t *testing.T) {
	actual1 := cloudevents.NewEvent()
	actual2 := cloudevents.NewEvent()
	actual3 := cloudevents.NewEvent()
	instance := TestAggregate{t: t}
	instance.SetID("id")
	// append an event and expect that the count is 1
	instance.AppendEvent(&actual1)
	assert.Equal(t, 1, len(instance.UncommittedEvents()))
	assert.Equal(t, &actual1, instance.uncommittedEvents[0])
	// append another two events and expect that the count is 3
	instance.AppendEvent(&actual2, &actual3)
	assert.Equal(t, 3, len(instance.UncommittedEvents()))
	assert.Equal(t, &actual2, instance.uncommittedEvents[1])
	assert.Equal(t, &actual3, instance.uncommittedEvents[2])
	// clear the events and expect that it clears the uncommitted events queue
	instance.ClearEvents()
	assert.Equal(t, 0, len(instance.UncommittedEvents()))
}
