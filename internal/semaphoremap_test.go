package internal

import (
	"context"
	"testing"

	helper "github.com/loadimpact/resolvent/resolventtest"
	"github.com/stretchr/testify/require"
)

func TestSemaphoreMap_Procure(t *testing.T) {
	t.Parallel()
	t.Run("singular_1", func(t *testing.T) {
		semaphores := NewSemaphoreMap(1)
		err := semaphores.Procure(context.Background(), "a")
		require.NoError(t, err, "procure failed")
	})
	t.Run("singular_3", func(t *testing.T) {
		semaphores := NewSemaphoreMap(3)
		for i := range helper.Three() {
			err := semaphores.Procure(context.Background(), "a")
			require.NoErrorf(t, err, "procure %d failed", i)
		}
	})
	t.Run("plural_1", func(t *testing.T) {
		semaphores := NewSemaphoreMap(1)
		for _, key := range helper.ABC() {
			err := semaphores.Procure(context.Background(), key)
			require.NoErrorf(t, err, "procure %s failed", key)
		}
	})
	t.Run("plural_3", func(t *testing.T) {
		semaphores := NewSemaphoreMap(3)
		for _, key := range helper.ABC() {
			for i := range helper.Three() {
				err := semaphores.Procure(context.Background(), key)
				require.NoErrorf(t, err, "procure %s %d failed", key, i)
			}
		}
	})
}

func TestSemaphoreMap_Vacate(t *testing.T) {
	t.Parallel()
	t.Run("invalid", func(t *testing.T) {
		semaphores := NewSemaphoreMap(1)
		require.Panics(t, func() {
			semaphores.Vacate("invalid-key")
		})
	})
	t.Run("singular_1", func(t *testing.T) {
		semaphores := NewSemaphoreMap(1)
		err := semaphores.Procure(context.Background(), "a")
		require.NoError(t, err, "procure failed")
		semaphores.Vacate("a")
	})
	t.Run("singular_3", func(t *testing.T) {
		semaphores := NewSemaphoreMap(3)
		for i := range helper.Three() {
			err := semaphores.Procure(context.Background(), "a")
			require.NoErrorf(t, err, "procure %d failed", i)
		}
		for range helper.Three() {
			semaphores.Vacate("a")
		}
	})
	t.Run("multiple", func(t *testing.T) {
		semaphores := NewSemaphoreMap(1)
		for _, key := range helper.ABC() {
			err := semaphores.Procure(context.Background(), key)
			require.NoErrorf(t, err, "procure %s failed", key)
		}
		for _, key := range helper.ABC() {
			semaphores.Vacate(key)
		}
	})
}

func TestSemaphoreMap_Reuse(t *testing.T) {
	t.Parallel()
	t.Run("singular_1", func(t *testing.T) {
		semaphores := NewSemaphoreMap(1)
		for i := range helper.Two() {
			err := semaphores.Procure(context.Background(), "a")
			require.NoErrorf(t, err, "procure %d failed", i)
			semaphores.Vacate("a")
		}
	})
	t.Run("singular_3", func(t *testing.T) {
		semaphores := NewSemaphoreMap(1)
		for i := range helper.Four() {
			err := semaphores.Procure(context.Background(), "a")
			require.NoErrorf(t, err, "procure %d failed", i)
			semaphores.Vacate("a")
		}
	})
	t.Run("plural_1", func(t *testing.T) {
		semaphores := NewSemaphoreMap(1)
		for _, key := range helper.ABC() {
			for i := range helper.Two() {
				err := semaphores.Procure(context.Background(), key)
				require.NoErrorf(t, err, "procure %s %d failed", key, i)
				semaphores.Vacate(key)
			}
		}
	})
	t.Run("plural_3", func(t *testing.T) {
		semaphores := NewSemaphoreMap(1)
		for _, key := range helper.ABC() {
			for i := range helper.Four() {
				err := semaphores.Procure(context.Background(), key)
				require.NoErrorf(t, err, "procure %s %d failed", key, i)
				semaphores.Vacate(key)
			}
		}
	})
}

func TestSemaphoreMap_Await(t *testing.T) {
	t.Parallel()
	t.Run("singular_1", func(t *testing.T) {
		semaphores := NewSemaphoreMap(1)
		err := semaphores.Procure(context.Background(), "a")
		require.NoError(t, err, "preliminary procure failed")
		procured := make(chan struct{})
		go func() {
			err := semaphores.Procure(context.Background(), "a")
			require.NoError(t, err, "awaiting procure failed")
			close(procured)
		}()
		<-semaphores.semaphore["a"].(*semaphore).procuring
		semaphores.Vacate("a")
		<-procured
	})
	t.Run("singular_3", func(t *testing.T) {
		semaphores := NewSemaphoreMap(1)
		err := semaphores.Procure(context.Background(), "a")
		require.NoError(t, err, "preliminary procure failed")
		procured := make(chan struct{})
		for i := range helper.Three() {
			go func(i int) {
				err := semaphores.Procure(context.Background(), "a")
				require.NoErrorf(t, err, "awaiting procure %d failed", i)
				procured <- struct{}{}
			}(i)
		}
		for range helper.Three() {
			<-semaphores.semaphore["a"].(*semaphore).procuring
			semaphores.Vacate("a")
			<-procured
		}
	})
	t.Run("plural_1", func(t *testing.T) {
		semaphores := NewSemaphoreMap(1)
		procured := make(map[string]chan struct{})
		for _, key := range helper.ABC() {
			err := semaphores.Procure(context.Background(), key)
			require.NoErrorf(t, err, "preliminary procure %s failed", key)
			procured[key] = make(chan struct{})
		}
		for _, key := range helper.ABC() {
			go func(key string) {
				err := semaphores.Procure(context.Background(), key)
				require.NoErrorf(t, err, "awaiting procure %s failed", key)
				close(procured[key])
			}(key)
		}
		for _, key := range helper.ABC() {
			<-semaphores.semaphore[key].(*semaphore).procuring
			semaphores.Vacate(key)
			<-procured[key]
		}
	})
	t.Run("plural_3", func(t *testing.T) {
		semaphores := NewSemaphoreMap(1)
		procured := make(map[string]chan struct{})
		for _, key := range helper.ABC() {
			err := semaphores.Procure(context.Background(), key)
			require.NoErrorf(t, err, "preliminary procure %s failed", key)
			procured[key] = make(chan struct{})
		}
		for _, key := range helper.ABC() {
			for i := range helper.Three() {
				go func(key string, i int) {
					err := semaphores.Procure(context.Background(), key)
					message := "awaiting procure %s %d failed"
					require.NoErrorf(t, err, message, key, i)
					procured[key] <- struct{}{}
				}(key, i)
			}
		}
		for _, key := range helper.ABC() {
			for range helper.Three() {
				<-semaphores.semaphore[key].(*semaphore).procuring
				semaphores.Vacate(key)
				<-procured[key]
			}
		}
	})
}

func TestSemaphoreMap_Abandon(t *testing.T) {
	t.Parallel()
	semaphores := NewSemaphoreMap(1)
	err := semaphores.Procure(context.Background(), "a")
	require.NoError(t, err, "preliminary procure failed")
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err = semaphores.Procure(ctx, "a")
	require.EqualError(t, err, "context canceled")
}
