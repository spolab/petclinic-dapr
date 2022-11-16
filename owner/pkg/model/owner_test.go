package model_test

import (
	"context"
	"testing"

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
