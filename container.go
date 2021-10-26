package di

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

var (
	ErrServiceNotRegistered = errors.New("service not registered")
)

type ServiceName string

//type Service interface {
//	//ServiceName() ServiceName
//}

type ServiceFactory func(ctx context.Context) interface{}

type registrationType string

const (
	registrationTypeTransient = "transient"
	registrationTypeSingleton = "singleton"
	registrationTypeScoped = "scoped"
)

type serviceRegistration struct {
	registrationType registrationType
	serviceName ServiceName
	serviceFactory ServiceFactory
	serviceInstance interface{}
}

type Container interface {
	RegisterTransient(name ServiceName, factory ServiceFactory) error
	RegisterSingleton(name ServiceName, factory ServiceFactory) error
	RegisterScoped(name ServiceName, factory ServiceFactory) error
	ResolveService(ctx context.Context, name ServiceName) (interface{}, error)
}

type container struct {
	ctx              context.Context
	ctxID            ContextID
	registrations    map[ServiceName]serviceRegistration
	registrationsMu sync.RWMutex
	scopedServices   map[ServiceName]map[ContextID]interface{}
	scopedServicesMu sync.Mutex
}

func (c *container) RegisterTransient(name ServiceName, sf ServiceFactory) error {
	c.register(serviceRegistration{
		registrationType: registrationTypeTransient,
		serviceName:      name,
		serviceFactory:   sf,
	})
	return nil
}

func (c *container) RegisterSingleton(name ServiceName, sf ServiceFactory) error {
	c.register(serviceRegistration{
		registrationType: registrationTypeSingleton,
		serviceName:      name,
		serviceFactory:   sf,
		serviceInstance: sf(c.ctx),
	})
	return nil
}

func (c *container) RegisterScoped(name ServiceName, sf ServiceFactory) error {
	c.register(serviceRegistration{
		registrationType: registrationTypeScoped,
		serviceName:      name,
		serviceFactory:   sf,
	})
	return nil
}

func (c *container) ResolveService(ctx context.Context, name ServiceName) (interface{}, error) {
	c.registrationsMu.RLock()
	defer c.registrationsMu.RUnlock()

	sr, ok := c.registrations[name]
	if !ok {
		return nil, ErrServiceNotRegistered
	}
	return c.resolveService(ctx, sr)
}

func (c *container) register(sr serviceRegistration) {
	c.registrationsMu.Lock()
	defer c.registrationsMu.Unlock()
	c.registrations[sr.serviceName] = sr
}

func (c *container) resolveService(ctx context.Context, sr serviceRegistration) (interface{}, error) {
	switch rt := sr.registrationType; rt {
	case registrationTypeSingleton:
		return c.resolveSingletonService(ctx, sr)
	case registrationTypeTransient:
		return c.resolveTransientService(ctx, sr)
	case registrationTypeScoped:
		return c.resolveScopedService(ctx, sr)
	default:
		return nil, fmt.Errorf(`unknown registration type: "%s"`, rt)
	}
}

func (c *container) resolveSingletonService(_ context.Context, sr serviceRegistration) (interface{}, error) {
	return sr.serviceInstance, nil
}

func (c *container) resolveTransientService(ctx context.Context, sr serviceRegistration) (interface{}, error) {
	return sr.serviceFactory(ctx), nil
}

func (c *container) resolveScopedService(ctx context.Context, sr serviceRegistration) (interface{}, error) {
	id, err := GetContextID(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get context ID: %w", err)
	}

	c.scopedServicesMu.Lock()
	defer c.scopedServicesMu.Unlock()

	serviceContexts, ok := c.scopedServices[sr.serviceName]
	if !ok {
		serviceContexts = map[ContextID]interface{}{}
		c.scopedServices[sr.serviceName] = serviceContexts
	}

	service, ok := serviceContexts[id]
	if !ok {
		service = sr.serviceFactory(ctx)
		serviceContexts[id] = service
	}
	return service, nil
}

func NewContainer(ctx context.Context) (Container, error) {
	ctxID, err := GetContextID(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get context ID: %w", err)
	}

	return &container{
		ctx:              ctx,
		ctxID:            ctxID,
		registrations:    map[ServiceName]serviceRegistration{},
		registrationsMu:  sync.RWMutex{},
		scopedServices: map[ServiceName]map[ContextID]interface{}{},
		scopedServicesMu: sync.Mutex{},
	}, nil
}
