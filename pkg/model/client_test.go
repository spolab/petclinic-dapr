package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Create a test with the factory method and check everything is as expected
func TestClientActorFactory(t *testing.T) {
	instance, ok := ClientActorFactory().(*ClientActor)
	require.True(t, ok)
	assert.NotNil(t, instance.UncommittedEvents())
	assert.NotNil(t, instance.Lifecycle)
}

func TestRegister(t *testing.T) {

}
