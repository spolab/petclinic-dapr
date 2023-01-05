package framework

import cloudevents "github.com/cloudevents/sdk-go/v2"

type EventSourcedState interface {
	// EventApplicator is responsible for applying a stream of events to a state of type S
	Apply(events ...*cloudevents.Event) error
	// Check the aggregate invariantes
	Check() error
}
