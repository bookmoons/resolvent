package internal

import (
	"context"
)

// Semaphore is a channel based semaphore.
type Semaphore interface {
	Procure(ctx context.Context) error
	Vacate()
}

type semaphore struct {
	active    chan struct{}
	procuring chan struct{}
}

// NewSemaphore returns a channel based semaphore.
func NewSemaphore(capacity uint16) *semaphore {
	return &semaphore{
		active:    make(chan struct{}, capacity),
		procuring: make(chan struct{}),
	}
}

// Procure allocates 1 unit of capacity when available.
// If the context closes, procurement is abandoned and the error is returned.
// Exactly 1 vacate must eventually follow each successful procurement.
func (s *semaphore) Procure(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case s.procuring <- struct{}{}:
	case s.active <- struct{}{}:
		return nil
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	case s.active <- struct{}{}:
		return nil
	}
}

// Vacate releases 1 unit of capacity.
func (s *semaphore) Vacate() {
	select {
	case <-s.active:
	default:
		panic("vacated unprocured semaphore")
	}
}
