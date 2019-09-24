package address

import (
	"context"
	"testing"
	"time"

	"github.com/loadimpact/resolvent"
	"github.com/loadimpact/resolvent/double"
	helper "github.com/loadimpact/resolvent/resolventtest"
	"github.com/miekg/dns"
	"github.com/stretchr/testify/require"
)

func TestQuery(t *testing.T) {
	t.Parallel()
	t.Run("invalid max", func(t *testing.T) {
		stub := double.NewStubQuerier([]*dns.Msg{})
		_, err := New(stub, 0)
		require.EqualError(t, err, "invalid max (0): must be positive")
	})
	t.Run("1 response", func(t *testing.T) {
		messages := helper.MakeMessages(t, []string{
			"year.test A 198.51.100.1",
		})
		stub := double.NewStubQuerier(messages)
		querier := construct(t, stub, 10)
		result := helper.SimpleQuery(t, querier, "query, failed")
		helper.DeepEqual(t, messages[0], result, "incorrect response")
	})
	t.Run("3 responses", func(t *testing.T) {
		messages := helper.MakeMessages(t, []string{
			"month.test A 198.51.100.2",
			"week.test A 198.51.100.3",
			"day.test A 198.51.100.4",
		})
		stub := double.NewStubQuerier(messages)
		querier := construct(t, stub, 10)
		results := []*dns.Msg{
			helper.SimpleQuery(t, querier, "query 1 failed"),
			helper.SimpleQuery(t, querier, "query 2 failed"),
			helper.SimpleQuery(t, querier, "query 3 failed"),
		}
		for i, message := range messages {
			result := results[i]
			helper.DeepEqual(t, message, result, "incorrect response")
		}
	})
	t.Run("respect max", func(t *testing.T) {
		stall := double.NewStallQuerier()
		querier := construct(t, stall, 2)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		go func() {
			helper.CancelableQuery(ctx, querier)
		}()
		go func() {
			helper.CancelableQuery(ctx, querier)
		}()
		<-stall.Received
		<-stall.Received
		_, err := helper.TimedQuery(querier, 100*time.Millisecond)
		require.EqualError(t, err, "context deadline exceeded")
	})
}

func construct(
	t *testing.T,
	underlying resolvent.Querier,
	max uint16,
) *addressLimitingQuerier {
	querier, err := New(underlying, max)
	require.NoError(t, err, "construct failed")
	return querier
}
