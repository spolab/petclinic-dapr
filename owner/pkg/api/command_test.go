package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	cmd := RegisterOwner{
		Id:         "id",
		Salutation: "salutation",
		Surname:    "surname",
		Name:       "name",
		Phone:      "phone",
		Email:      "email@google.com",
	}
	assert.Nil(t, cmd.Validate())
}
