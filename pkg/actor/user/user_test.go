package user_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/spolab/petstore/pkg/actor/user"
	"github.com/spolab/petstore/pkg/mock_client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserActor(t *testing.T) {
	ctrl := gomock.NewController(t)
	client := mock_client.NewMockClient(ctrl)
	factory := user.UserActorFactory(client, "broker", "topic")
	actual := factory()
	require.NotNil(t, actual)
	assert.Equal(t, "user", actual.Type())
}
