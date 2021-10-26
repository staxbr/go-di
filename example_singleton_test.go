package di_test

import (
	"context"
	"fmt"
	"github.com/staxbr/go-di"
	"time"
)

func ExampleSingleton() {
	ctx := di.NewContext()
	c, err := di.NewContainer(ctx)
	if err != nil {
		panic(err)
	}

	err = c.RegisterSingleton(ClockServiceName, ClockServiceFactory)
	if err != nil {
		panic(err)
	}

	givenService1st, err := c.ResolveService(ctx, ClockServiceName)
	if err != nil {
		panic(err)
	}

	givenService2nd, err := c.ResolveService(ctx, ClockServiceName)
	if err != nil {
		panic(err)
	}

	otherCtx := di.ContextWithID(ctx)
	givenService3rd, err := c.ResolveService(otherCtx, ClockServiceName)
	if err != nil {
		panic(err)
	}

	// as a singleton, should be the same instance
	if givenService1st != givenService2nd {
		panic("mismatch")
	}
	if givenService2nd != givenService3rd {
		panic("mismatch")
	}

	fmt.Printf("%s\n", givenService1st.(*ClockService).Now())
	fmt.Printf("%s\n", givenService2nd.(*ClockService).Now())
	fmt.Printf("%s\n", givenService3rd.(*ClockService).Now())

	// Output: 2021-10-26 19:47:00 +0000 UTC
	// 2021-10-26 19:47:00 +0000 UTC
	// 2021-10-26 19:47:00 +0000 UTC
}

const ClockServiceName = di.ServiceName("Clock")

type ClockService struct {
	clock func() time.Time
}

func ClockServiceFactory(ctx context.Context) interface{} {
	return newFixedClock(time.Date(2021, time.October, 26, 19, 47, 0, 0, time.UTC))
}

func newFixedClock(fixedMoment time.Time) *ClockService {
	return &ClockService{
		clock: func() time.Time {
			return fixedMoment
		},
	}
}

func (s ClockService) Now() time.Time {
	return s.clock()

}
