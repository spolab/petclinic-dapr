package model

import (
	"context"
	"encoding/json"

	"github.com/dapr/go-sdk/client"
)

type Repository struct {
	dapr  client.Client
	appId string
	topic string
}

func (r *Repository) GetById(ctx context.Context, id string) (*Owner, error) {
	state, err := r.dapr.GetState(ctx, r.appId, id, nil)
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
func (r *Repository) Save(ctx context.Context, owner *Owner) error {
	bytes, err := json.Marshal(owner.State)
	if err != nil {
		return err
	}
	if err := r.dapr.SaveState(ctx, r.appId, "owner-state", bytes, nil); err != nil {
		return err
	}
	for _, event := range owner.Events {
		if err := r.dapr.PublishEvent(ctx, "owner-pubsub", r.topic, event); err != nil {
			// TODO how do we tell which events have not been published?
			return err
		}
	}
	return nil
}

func NewRepository(c client.Client, id string) (*Repository, error) {
	return &Repository{
		dapr:  c,
		appId: id,
	}, nil
}
