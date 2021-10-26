# go-di

## Usage

### Singleton

```go
package main

import (
	"context"
	"github.com/staxbr/go-di"
	"time"
)

func main() {
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
}

const ClockServiceName = di.ServiceName("Clock")

type ClockService struct {
	clock func() time.Time
}

func ClockServiceFactory(ctx context.Context) interface{} {
	return newSystemClock()
}

func newSystemClock() ClockService {
	return ClockService{
		clock: time.Now,
    }
}

func (s ClockService) Now() time.Time {
    return s.clock()
}
```