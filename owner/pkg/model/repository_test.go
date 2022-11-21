package model_test

import (
	"context"
	"testing"

	"github.com/dapr/go-sdk/client"
	"github.com/spolab/petclinic/owner/pkg/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Tests the save of an aggregate root with one event
func TestSaveHappyPath(t *testing.T) {
	c, err := client.NewClient()
	require.Nil(t, err)
	repo, err := model.NewRepository(c, "owner-state", "owner-pubsub", "owner")
	require.Nil(t, err)
	require.NotNil(t, repo)
	aggregate := model.Owner{
		Id: "id",
	}
	assert.Nil(t, repo.Save(context.TODO(), &aggregate))
}
