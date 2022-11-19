package model_test

import (
	"context"
	"testing"

	ce "github.com/cloudevents/sdk-go/v2"
	"github.com/spolab/petclinic/owner/pkg/api"
	"github.com/spolab/petclinic/owner/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestApplyOwnerRegistered(t *testing.T) {
	instance := model.Owner{}
	instance.Apply(context.TODO(), api.OwnerRegistered{Id: "id", Salutation: "salutation", Surname: "surname", Name: "name", Phone: "phone", Email: "foo@bar.com"})
	assert.Equal(t, "id", instance.Id)
	assert.Equal(t, "salutation", instance.State.Salutation)
	assert.Equal(t, "surname", instance.State.Surname)
	assert.Equal(t, "name", instance.State.Name)
	assert.Equal(t, "phone", instance.State.Phone)
	assert.Equal(t, "foo@bar.com", instance.State.Email)
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
		Phone:      "phone",
		Email:      "foobar@baz.com",
	}
	instance.Register(context.TODO(), cmd)
	// Inspect the contents of the aggregate root
	assert.Equal(t, "id", instance.Id)
	assert.Equal(t, "salutation", instance.State.Salutation)
	assert.Equal(t, "surname", instance.State.Surname)
	assert.Equal(t, "name", instance.State.Name)
	assert.Equal(t, 1, len(instance.UncommittedEvents))
	// Inspect the contents of the event
	event, ok := instance.UncommittedEvents[0].(api.OwnerRegistered)
	assert.True(t, ok)
	assert.Equal(t, "id", event.Id)
	assert.Equal(t, "salutation", event.Salutation)
	assert.Equal(t, "surname", event.Surname)
	assert.Equal(t, "name", event.Name)
}

func TestRegisterWithNilCommandReturnsError(t *testing.T) {
	instance := model.Owner{}
	assert.NotNil(t, instance.Register(context.TODO(), nil))
}
