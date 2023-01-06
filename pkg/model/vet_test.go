package model

import (
	"context"
	"testing"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	"github.com/spolab/petstore/gen/mock/dapr/actor"
	"github.com/spolab/petstore/pkg/command"
	"github.com/spolab/petstore/pkg/event"
	"github.com/spolab/petstore/pkg/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Tests that hydrating a Vet from an event stream gives us what we expect
func TestVetLoadOk(t *testing.T) {
	vetId := "1234"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sm := actor.NewMockStateManager(ctrl)
	//
	// Creates the stream of events that the actor will load
	//
	registered := cloudevents.NewEvent()
	registered.SetType(event.TypeVetRegisteredV1)
	registered.SetData(cloudevents.ApplicationJSON, event.VetRegistered{Id: vetId, Name: "name", Surname: "surname", Phone: "phone", Email: "mail@mail.com"})
	stream := []*cloudevents.Event{
		&registered,
	}
	//
	// Creates the actor instance we are goint to test
	//
	instance := Vet{}
	instance.SetID(vetId)
	instance.SetStateManager(sm)
	instance.events = stream
	//
	// Mock the state manager to return the stream
	//
	sm.EXPECT().Contains(KeyEvents).Return(true, nil)
	sm.EXPECT().Get(KeyEvents, &instance.events).Return(nil)
	//
	// Load the aggregate from the events
	//
	require.Nil(t, instance.Load())
	state := instance.state
	assert.Equal(t, vetId, state.Id)
	assert.Equal(t, "name", state.Name)
	assert.Equal(t, "surname", state.Surname)
	assert.Equal(t, "phone", state.Phone)
	assert.Equal(t, "mail@mail.com", state.Email)
	assert.Equal(t, 0, state.Version)
	assert.False(t, state.Deleted)
	assert.Equal(t, 1, len(instance.events))
}

// Tests a valid vet registration
func TestRegisterVetOk(t *testing.T) {
	vetId := "1234"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sm := actor.NewMockStateManager(ctrl)
	//
	// Creates the actor instance without any prior events
	//
	instance := Vet{}
	instance.SetID(vetId)
	instance.SetStateManager(sm)
	instance.validate = validator.New()
	instance.events = []*cloudevents.Event{}
	//
	// Expected array of events
	//
	vr := &event.VetRegistered{Id: vetId, Name: "name", Surname: "surname", Phone: "phone", Email: "mail@mail.com"}
	registered := event.CloudEvent(event.FromSource("vet"), event.OfType(event.TypeVetRegisteredV1), event.WithDataAsJSON(vr))
	//
	// Execute the request
	//
	sm.EXPECT().Contains(KeyEvents).Return(false, nil)
	sm.EXPECT().Set(KeyEvents, MatchesEvents(registered)).Return(nil)
	_, err := instance.Register(context.Background(), &command.RegisterVetCommand{Name: "name", Surname: "surname", Phone: "phone", Email: "mail@mail.com"})
	require.Nil(t, err)
	//
	// Verify that the state matches the command
	//
	state := instance.state
	assert.Equal(t, vetId, state.Id)
	assert.Equal(t, "name", state.Name)
	assert.Equal(t, "surname", state.Surname)
	assert.Equal(t, "phone", state.Phone)
	assert.Equal(t, "mail@mail.com", state.Email)
	assert.Equal(t, 0, state.Version)
	assert.False(t, state.Deleted)
	assert.Equal(t, 1, len(instance.events))
	//
	// Verify that the events queue and the return type contain the right event
	//
	assert.Equal(t, 1, len(instance.events))
}

func MatchesEvents(events ...*cloudevents.Event) gomock.Matcher {
	return test.CloudEventArrayMatcher{Expected: events}
}
