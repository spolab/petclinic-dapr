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
package event

import (
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
)

func CloudEvent(options ...EventOption) *cloudevents.Event {
	result := cloudevents.NewEvent()
	for _, option := range options {
		option.Apply(&result)
	}
	result.SetID(uuid.NewString())
	result.SetTime(time.Now())
	return &result
}

type EventOption interface {
	Apply(target *cloudevents.Event)
}

func WithDataAsJSON(data any) EventOption {
	return &WithDataAsOption{contentType: cloudevents.ApplicationJSON, data: data}
}

type WithDataAsOption struct {
	contentType string
	data        any
}

func (o *WithDataAsOption) Apply(event *cloudevents.Event) {
	event.SetData(o.contentType, o.data)
}

func FromSource(source string) EventOption {
	return &FromSourceOption{source: source}
}

type FromSourceOption struct {
	source string
}

func (o *FromSourceOption) Apply(event *cloudevents.Event) {
	event.SetSource(o.source)
}

func OfType(kind string) EventOption {
	return &OfTypeOption{kind: kind}
}

type OfTypeOption struct {
	kind string
}

func (o *OfTypeOption) Apply(event *cloudevents.Event) {
	event.SetType(o.kind)
}
