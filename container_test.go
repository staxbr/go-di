package di_test

import (
	"context"
	"github.com/staxbr/go-di"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestResolver(t *testing.T) {
	suite.Run(t, new(ResolverTest))
}

type ResolverTest struct {
	suite.Suite
}

func (t *ResolverTest) TestContainer_RegisterSingleton() {
	// Arrange
	ctx := di.NewContext()
	c, err := di.NewContainer(ctx)
	t.Require().NoError(err)

	myServiceName, myServiceFactory := t.defineNewService("MyService")
	err = c.RegisterSingleton(myServiceName, myServiceFactory)
	t.Require().NoError(err)

	// Act
	givenService1st, err := c.ResolveService(ctx, myServiceName)
	t.Require().NoError(err)
	givenService2nd, err := c.ResolveService(ctx, myServiceName)
	t.Require().NoError(err)
	givenService3rd, err := c.ResolveService(di.ContextWithID(ctx), myServiceName)
	t.Require().NoError(err)

	// Assert
	t.Same(givenService1st, givenService2nd)
	t.Same(givenService2nd, givenService3rd)
}

func (t *ResolverTest) TestContainer_RegisterScoped() {
	// Arrange
	ctx := di.NewContext()
	c, err := di.NewContainer(ctx)
	t.Require().NoError(err)

	// Act
	myServiceName, myServiceFactory := t.defineNewService("MyService")
	err = c.RegisterScoped(myServiceName, myServiceFactory)
	t.Require().NoError(err)

	givenService1st, err := c.ResolveService(ctx, myServiceName)
	t.Require().NoError(err)
	givenService2nd, err := c.ResolveService(ctx, myServiceName)
	t.Require().NoError(err)

	// Assert
	t.Same(givenService1st, givenService2nd)
}

func (t *ResolverTest) TestContainer_RegisterTransient() {
	// Arrange
	ctx := di.NewContext()
	c, err := di.NewContainer(ctx)
	t.Require().NoError(err)

	myServiceName, myServiceFactory := t.defineNewService("MyService")
	err = c.RegisterTransient(myServiceName, myServiceFactory)
	t.Require().NoError(err)

	// Act
	givenService1st, err := c.ResolveService(ctx, myServiceName)
	t.Require().NoError(err)
	givenService2nd, err := c.ResolveService(ctx, myServiceName)
	t.Require().NoError(err)

	// Assert
	t.NotSame(givenService1st, givenService2nd)
	t.NotNil(givenService1st,)
	t.NotNil(givenService2nd)
}

func (t *ResolverTest) TestResolveService_WhenServiceIsNotRegistered_ShouldReturnExpectedError() {
	ctx := di.NewContext()
	c, err := di.NewContainer(ctx)
	t.Require().NoError(err)

	givenService, err := c.ResolveService(ctx, "MyServiceName")

	t.Assert().Nil(givenService)
	t.Assert().ErrorIs(err, di.ErrServiceNotRegistered)
}

func (t *ResolverTest) defineNewService(name string) (di.ServiceName, di.ServiceFactory) {
	t.T().Helper()

	type serviceForResolverTest struct {
		name di.ServiceName // https://dave.cheney.net/2014/03/25/the-empty-struct
	}

	serviceName := di.ServiceName(name)
	serviceFactory := func(_ context.Context) interface{} {
		return &serviceForResolverTest{
			name: serviceName,
		}
	}
	return serviceName, serviceFactory
}

//func(s *serviceForResolverTest) set() {
//	s.v = time.Now()
//}