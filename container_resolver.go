package di

import (
	"context"
	"fmt"
	"sync"
)

var (
	containers = map[ContextID]Container{}
	containersMu = sync.Mutex{}
)

func ResolveContainer(ctx context.Context) Container {
	id, err := GetContextID(ctx)
	if err != nil {
		return newErrorContainer(fmt.Errorf("unable to get context ID: %w", err))
	}

	if _, ok := containers[id]; !ok {
		c, err := NewContainer(ctx)
		if err != nil {
			return newErrorContainer(fmt.Errorf("unable to create new container: %w", err))
		}
		containersMu.Lock()
		containers[id] = c
		containersMu.Unlock()
	}
	return containers[id]
}
