package model_test

import (
	"context"
	"testing"

	ce "github.com/cloudevents/sdk-go/v2"
	"github.com/spolab/petclinic/owner/pkg/api"
	"github.com/spolab/petclinic/owner/pkg/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
	require.Nil(t, instance.Register(context.TODO(), cmd))
	// Inspect the contents of the aggregate root
	assert.Equal(t, "id", instance.Id)
	assert.Equal(t, "salutation", instance.State.Salutation)
	assert.Equal(t, "surname", instance.State.Surname)
	assert.Equal(t, "name", instance.State.Name)
	assert.Equal(t, "phone", instance.State.Phone)
	assert.Equal(t, "foobar@baz.com", instance.State.Email)
	assert.Equal(t, 1, len(instance.UncommittedEvents))
	// Inspect the contents of the event
	event := instance.UncommittedEvents[0]
	var data api.OwnerRegistered
	require.Nil(t, event.DataAs(&data))
	assert.Equal(t, "id", data.Id)
	assert.Equal(t, "salutation", data.Salutation)
	assert.Equal(t, "surname", data.Surname)
	assert.Equal(t, "name", data.Name)
	assert.Equal(t, "phone", data.Phone)
	assert.Equal(t, "foobar@baz.com", data.Email)
}

func TestRegisterWithNilCommandReturnsError(t *testing.T) {
	instance := model.Owner{}
	assert.NotNil(t, instance.Register(context.TODO(), nil))
}
