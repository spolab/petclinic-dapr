package events

import (
	cloudevents "github.com/cloudevents/sdk-go/v2"
)

func CloudEvent(source string, kind string, data any) cloudevents.Event {
	result := cloudevents.NewEvent()
	result.SetSource(source)
	result.SetType(kind)
	result.SetData(cloudevents.ApplicationJSON, data)
	return result
}
