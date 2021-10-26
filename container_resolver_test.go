package di_test

import (
	"github.com/staxbr/go-di"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestResolveContainer_WhenNilContextIsGiven_ShouldResolveAnErrorContainer(t *testing.T) {
	t.Parallel()

	c := di.ResolveContainer(nil)
	require.NotNil(t, c)

	serviceName := di.ServiceName("MyServiceName")
	givenService, err := c.ResolveService(nil, serviceName)

	assert.Nil(t, givenService)
	assert.EqualError(t, err, "[error container] reason: unable to get context ID: ctx cannot be nil")
}

func TestResolveContainer_WhenTheSameContextIsGiven_ShouldResolveTheSameContainerInstance(t *testing.T) {
	t.Parallel()

	ctx := di.NewContext()
	resolvedContainer1 := di.ResolveContainer(ctx)
	require.NotNil(t, resolvedContainer1)
	resolvedContainer2 := di.ResolveContainer(ctx)
	require.NotNil(t, resolvedContainer2)

	assert.Same(t, resolvedContainer1, resolvedContainer2)
}
