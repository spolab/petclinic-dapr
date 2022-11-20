package model

import (
	"context"
	"encoding/json"
	"reflect"

	ce "github.com/cloudevents/sdk-go/v2"
	"github.com/dapr/go-sdk/client"
	"github.com/spolab/petclinic/owner/pkg/api"
)

type Repository struct {
	dapr        client.Client
	storeName   string
	brokerName  string
	brokerTopic string
}

func (r *Repository) GetById(ctx context.Context, id string) (*Owner, error) {
	state, err := r.dapr.GetState(ctx, r.storeName, id, nil)
	if err != nil {
		return nil, err
	}
	result := Owner{Id: id, Version: state.Etag}
	if err := json.Unmarshal(state.Value, &result.State); err != nil {
		return nil, err
	}
	return &result, nil
}

// Saves the state of the aggregate and emits all the new events
func (r *Repository) Save(ctx context.Context, owner *Owner, so ...client.StateOption) error {
	bytes, err := json.Marshal(owner.State)
	if err != nil {
		return err
	}
	if err := r.dapr.SaveState(ctx, r.storeName, owner.Id, bytes, nil, so...); err != nil {
		return err
	}
	meta := make(map[string]string)
	for _, event := range owner.UncommittedEvents {
		meta[api.MetaEventType] = reflect.TypeOf(event).Name()
		if err := r.dapr.PublishEvent(ctx, r.brokerName, r.brokerTopic, event, client.PublishEventWithContentType(ce.ApplicationCloudEventsJSON)); err != nil {
			// TODO how do we tell which events have not been published?
			return err
		}
	}
	// Empty out the list of uncommitted events and exit
	owner.UncommittedEvents = []ce.Event{}
	return nil
}

func NewRepository(c client.Client, storeName string, brokerName string, brokerTopic string) (*Repository, error) {
	return &Repository{
		dapr:        c,
		storeName:   storeName,
		brokerName:  brokerName,
		brokerTopic: brokerTopic,
	}, nil
}
