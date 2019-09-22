package internal

import (
	"context"

	"github.com/pkg/errors"
)

// SemaphoreMap is a string keyed map of semaphores.
type SemaphoreMap interface {
	Procure(ctx context.Context, key string) error
	Vacate(key string)
}

type semaphoreMapEntry struct {
	semaphore Semaphore
	active    uint16
}

type semaphoreMap struct {
	lock      Semaphore
	semaphore map[string]Semaphore
	active    map[string]uint16
	capacity  uint16
}

// NewSemaphoreMap returns a map of semaphores.
func NewSemaphoreMap(capacity uint16) *semaphoreMap {
	return &semaphoreMap{
		lock:      NewSemaphore(1),
		semaphore: make(map[string]Semaphore),
		active:    make(map[string]uint16),
		capacity:  capacity,
	}
}

// Procure allocates 1 unit of the specified key's capacity when available.
// If the context closes, procurement is abandoned and the error is returned.
// Each successful procurement of a key must eventually be followed by
// exactly 1 vacation of the same key.
func (m *semaphoreMap) Procure(ctx context.Context, key string) (err error) {
	semaphore, err := m.acquireEntry(ctx, key)
	if err != nil {
		return
	}
	err = semaphore.Procure(ctx)
	if err != nil {
		m.abandonEntry(key)
		return
	}
	return
}

// Vacate releases 1 unit of the specified key's capacity.
func (m *semaphoreMap) Vacate(key string) {
	m.releaseEntry(key)
}

func (m *semaphoreMap) acquireEntry(
	ctx context.Context,
	key string,
) (semaphore Semaphore, err error) {
	if err = m.lock.Procure(ctx); err != nil {
		return
	}
	defer m.lock.Vacate()
	if _, ok := m.semaphore[key]; !ok {
		m.semaphore[key] = NewSemaphore(m.capacity)
	}
	if m.active[key] == ^uint16(0) {
		// Active count overflow
		err = errors.New("exceeded maximum waiting procurements")
		return
	}
	m.active[key]++
	return m.semaphore[key], nil
}

func (m *semaphoreMap) abandonEntry(key string) {
	m.lock.Procure(context.Background())
	defer m.lock.Vacate()
	if _, ok := m.semaphore[key]; !ok {
		panic("invalid key")
	}
	if m.active[key] == 0 {
		// Active count underflow
		panic("abandoned unprocured entry")
	}
	m.active[key]--
	if m.active[key] == 0 {
		delete(m.semaphore, key)
		delete(m.active, key)
	}
}

func (m *semaphoreMap) releaseEntry(key string) {
	m.lock.Procure(context.Background())
	defer m.lock.Vacate()
	if _, ok := m.semaphore[key]; !ok {
		panic("invalid key")
	}
	if m.active[key] == 0 {
		// Active count underflow
		panic("released unprocured entry")
	}
	m.semaphore[key].Vacate()
	m.active[key]--
	if m.active[key] == 0 {
		delete(m.semaphore, key)
		delete(m.active, key)
	}
}
