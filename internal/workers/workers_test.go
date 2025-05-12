package workers_test

import (
	"context"
	"testing"
	"time"

	"github.com/henrywhitaker3/boiler"
	"github.com/henrywhitaker3/shorturl/internal/test"
	"github.com/henrywhitaker3/shorturl/internal/workers"
	"github.com/redis/rueidis"
	"github.com/stretchr/testify/require"
)

type testWorker struct {
	interval   time.Duration
	timeout    time.Duration
	executions int
}

func (t *testWorker) Name() string {
	return "tester"
}

func (t *testWorker) Interval() workers.Interval {
	return workers.NewInterval(t.interval)
}

func (t *testWorker) Timeout() time.Duration {
	return t.timeout
}

func (t *testWorker) Run(ctx context.Context) error {
	t.executions++
	return nil
}

func (t *testWorker) Executions() int {
	return t.executions
}

func TestItRunsWorkers(t *testing.T) {
	b := test.Boiler(t)
	redis, err := boiler.Resolve[rueidis.Client](b)
	require.Nil(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	worker := &testWorker{
		interval:   time.Millisecond * 250,
		timeout:    time.Second,
		executions: 0,
	}

	runner, err := workers.NewRunner(ctx, redis)
	require.Nil(t, err)
	require.Nil(t, runner.Register(worker))
	runner.Run()
	defer runner.Stop()
	time.Sleep(time.Millisecond * 500)

	require.Equal(t, 1, worker.Executions())
}
