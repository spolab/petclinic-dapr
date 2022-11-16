package model

import ce "github.com/cloudevents/sdk-go/v2"

type Aggregate[T any] struct {
	Id                string
	Version           string
	UncommittedEvents []ce.Event
	State             T
}
