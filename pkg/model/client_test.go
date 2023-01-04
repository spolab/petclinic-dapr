package model

import (
	"context"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	"github.com/spolab/petstore/gen/mock/dapr/actor"
	"github.com/spolab/petstore/pkg/command"
	"github.com/stretchr/testify/require"
)

var validRegisterCommand = &command.RegisterClientCommand{Salutation: "salutation", Name: "name", Surname: "surname", Phone: "phone", Email: "mail@mail.com"}

func TestRegister(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sm := actor.NewMockStateManager(ctrl)
	instance := Client{validate: validator.New()}
	instance.SetStateManager(sm)
	instance.SetID("id")
	// pretend to not exist
	sm.EXPECT().Contains(KeyEvents).Return(false, nil)
	// execute test
	_, err := instance.Register(context.TODO(), validRegisterCommand)
	require.Nil(t, err)
}

func TestLoad(t *testing.T) {

}
