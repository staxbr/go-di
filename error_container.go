package di

import (
	"context"
	"fmt"
)

type errorContainer struct {
	err error
}

func newErrorContainer(err error) *errorContainer {
	return &errorContainer{err: err}
}

func (c *errorContainer) RegisterTransient(_ ServiceName, _ ServiceFactory) error {
	return c.err
}

func (c *errorContainer) RegisterSingleton(_ ServiceName, _ ServiceFactory) error {
	return c.err
}

func (c *errorContainer) RegisterScoped(_ ServiceName, _ ServiceFactory) error {
	return c.err
}

func (c *errorContainer) ResolveService(_ context.Context, _ ServiceName) (interface{}, error) {
	return nil, fmt.Errorf("[error container] reason: %w", c.err)
}

