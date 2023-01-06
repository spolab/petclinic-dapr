package framework

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEventSourcedCommandLifecycle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := NewMockEventSourcedRepository(ctrl)
	aggregate := &TestAggregate{t: t}
	aggregate.SetID("id")
	instance := &EventSourcedCommandLifecycle{Repository: repo}
	repo.EXPECT().Load(aggregate).Return(nil)
	repo.EXPECT().Save(aggregate).Return(nil)
	err := instance.Execute(aggregate, func() error {
		aggregate.ExecuteCount++
		return nil
	})
	require.Nil(t, err)
	assert.Equal(t, 1, aggregate.ApplyCount)
	assert.Equal(t, 2, aggregate.CheckCount) // 2 invocations, pre- and post-execution
	assert.Equal(t, 1, aggregate.ExecuteCount)
}
