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
	require.NotNil(t, instance)
	assert.NotNil(t, instance.commandLifecycle)
	assert.NotNil(t, instance.mutator)
	assert.NotNil(t, instance.snapshot)
	assert.NotNil(t, instance.uncommittedEvents)
}
