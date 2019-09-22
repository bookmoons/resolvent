package internal

import (
	"context"
	"testing"

	helper "github.com/loadimpact/resolvent/resolventtest"
	"github.com/stretchr/testify/require"
)

func TestSemaphore_Procure(t *testing.T) {
	t.Parallel()
	t.Run("1", func(t *testing.T) {
		semaphore := NewSemaphore(1)
		err := semaphore.Procure(context.Background())
		require.NoError(t, err, "procure failed")
	})
	t.Run("3", func(t *testing.T) {
		semaphore := NewSemaphore(3)
		for i := range helper.Three() {
			err := semaphore.Procure(context.Background())
			require.NoErrorf(t, err, "procure %d failed", i)
		}
	})
}

func TestSemaphore_Vacate(t *testing.T) {
	t.Parallel()
	t.Run("invalid", func(t *testing.T) {
		semaphore := NewSemaphore(1)
		require.Panics(t, func() {
			semaphore.Vacate()
		})
	})
	t.Run("1", func(t *testing.T) {
		semaphore := NewSemaphore(1)
		err := semaphore.Procure(context.Background())
		require.NoError(t, err, "procure failed")
		semaphore.Vacate()
	})
	t.Run("3", func(t *testing.T) {
		semaphore := NewSemaphore(3)
		for i := range helper.Three() {
			err := semaphore.Procure(context.Background())
			require.NoErrorf(t, err, "procure %d failed", i)
		}
		for range helper.Three() {
			semaphore.Vacate()
		}
	})
}

func TestSemaphore_Reuse(t *testing.T) {
	t.Parallel()
	t.Run("1", func(t *testing.T) {
		semaphore := NewSemaphore(1)
		for i := range helper.Two() {
			err := semaphore.Procure(context.Background())
			require.NoErrorf(t, err, "procure %d failed", i)
			semaphore.Vacate()
		}
	})
	t.Run("3", func(t *testing.T) {
		semaphore := NewSemaphore(1)
		for i := range helper.Four() {
			err := semaphore.Procure(context.Background())
			require.NoErrorf(t, err, "procure %d failed", i)
			semaphore.Vacate()
		}
	})
}

func TestSemaphore_Await(t *testing.T) {
	t.Parallel()
	t.Run("1", func(t *testing.T) {
		semaphore := NewSemaphore(1)
		err := semaphore.Procure(context.Background())
		require.NoError(t, err, "preliminary procure failed")
		procured := make(chan struct{})
		go func() {
			err := semaphore.Procure(context.Background())
			require.NoError(t, err, "awaiting procure failed")
			close(procured)
		}()
		<-semaphore.procuring
		semaphore.Vacate()
		<-procured
	})
	t.Run("3", func(t *testing.T) {
		semaphore := NewSemaphore(1)
		err := semaphore.Procure(context.Background())
		require.NoError(t, err, "preliminary procure failed")
		procured := make(chan struct{})
		for i := range helper.Three() {
			go func(i int) {
				err := semaphore.Procure(context.Background())
				require.NoErrorf(t, err, "awaiting procure %d failed", i)
				procured <- struct{}{}
			}(i)
		}
		for range helper.Three() {
			<-semaphore.procuring
			semaphore.Vacate()
			<-procured
		}
	})
}

func TestSemaphore_Abandon(t *testing.T) {
	t.Parallel()
	semaphore := NewSemaphore(1)
	err := semaphore.Procure(context.Background())
	require.NoError(t, err, "preliminary procure failed")
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err = semaphore.Procure(ctx)
	require.EqualError(t, err, "context canceled")
}
