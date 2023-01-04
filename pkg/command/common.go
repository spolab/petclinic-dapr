package command

import cloudevents "github.com/cloudevents/sdk-go/v2"

const (
	StatusOK = iota
	StatusInvalid
	StatusError
)

type ActorResponse struct {
	Status  int                  `json:"status"`
	Message string               `json:"message,omitempty"`
	Events  []*cloudevents.Event `json:"events,omitempty"`
}
