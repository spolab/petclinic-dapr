package event

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCloudEvent(t *testing.T) {
	expData := map[string]string{"key": "value"}
	actual := CloudEvent(FromSource("source"), OfType("kind"), WithDataAsJSON(expData))
	require.NotNil(t, actual)
	assert.Equal(t, "1.0", actual.SpecVersion())
	_, err := uuid.Parse(actual.ID())
	assert.Nil(t, err)
	assert.NotNil(t, actual.Time())
	assert.Equal(t, "application/json", actual.DataMediaType())
	assert.Equal(t, "source", actual.Source())
	assert.Equal(t, "kind", actual.Type())
	var actData map[string]string
	require.Nil(t, actual.DataAs(&actData))
	assert.Equal(t, 1, len(actData))
	assert.Equal(t, "value", actData["key"])
}
