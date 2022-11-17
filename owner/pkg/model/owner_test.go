package model_test

import (
	"context"
	"encoding/json"
	"testing"

	ce "github.com/cloudevents/sdk-go/v2"
	"github.com/spolab/petclinic/owner/pkg/api"
	"github.com/spolab/petclinic/owner/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestApplyOwnerRegistered(t *testing.T) {
	instance := model.Owner{}
	instance.Apply(context.TODO(), model.NewOwnerCreated("id", "salutation", "surname", "name"))
	assert.Equal(t, "id", instance.Id)
	assert.Equal(t, "salutation", instance.State.Salutation)
	assert.Equal(t, "surname", instance.State.Surname)
	assert.Equal(t, "name", instance.State.Name)
}

func TestNewOwnerCreated(t *testing.T) {
	actual := model.NewOwnerCreated("id", "salutation", "surname", "name")
	assert.Equal(t, "1.0", actual.SpecVersion())
	assert.NotEmpty(t, actual.ID())
	assert.Equal(t, "application/json", actual.DataContentType())
}

func TestUnknownEventReturnsError(t *testing.T) {
	instance := model.Owner{}
	assert.NotNil(t, instance.Apply(context.TODO(), ce.NewEvent()))
}

func TestRegisterHappyPath(t *testing.T) {
	instance := model.Owner{}
	cmd := &api.RegisterOwner{
		Id:         "id",
		Salutation: "salutation",
		Surname:    "surname",
		Name:       "name",
	}
	instance.Register(context.TODO(), cmd)
	// Inspect the contents of the aggregate root
	assert.Equal(t, "id", instance.Id)
	assert.Equal(t, "salutation", instance.State.Salutation)
	assert.Equal(t, "surname", instance.State.Surname)
	assert.Equal(t, "name", instance.State.Name)
	assert.Equal(t, 1, len(instance.UncommittedEvents))
	// Inspect the contents of the event
	event := instance.UncommittedEvents[0]
	assert.NotNil(t, event.ID())
	assert.NotNil(t, "1.0", event.SpecVersion())
	assert.NotNil(t, model.OwnerRegisteredUrn, event.DataSchema())
	// Inspect the event payload
	var data api.OwnerRegistered
	assert.Nil(t, json.Unmarshal(event.Data(), &data))
	assert.Equal(t, "id", data.Id)
	assert.Equal(t, "salutation", data.Salutation)
	assert.Equal(t, "surname", data.Surname)
	assert.Equal(t, "name", data.Name)
}

func TestRegisterWithNilCommandReturnsError(t *testing.T) {
	instance := model.Owner{}
	assert.NotNil(t, instance.Register(context.TODO(), nil))
}
