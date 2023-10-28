package main

import (
	"context"
	"sync"
)

type muctx struct {
	mu   sync.Mutex
	chmu chan struct{}
}

// New returns new mutex with context
func New(mu ...sync.Mutex) *muctx {
	if len(mu) == 1 {
		return &muctx{
			mu:   mu[0],
			chmu: make(chan struct{}, 1),
		}
	}
	return &muctx{
		chmu: make(chan struct{}, 1),
	}
}

// Lock returns true if mutex set. False is context Done. Will lock process till one of this results
func (mx *muctx) Lock(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return false
	case mx.chmu <- struct{}{}:
		mx.mu.Lock()
		return true
	}
}

// Unlock release mutex
func (mx *muctx) Unlock() {
	<-mx.chmu
	mx.mu.Unlock()
}
