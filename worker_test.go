package workers

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestOneWorker(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool := NewWorkerPool(ctx, 1)

	job := func() {
		fmt.Println("run")
	}
	time.Sleep(2)
	{ // Submit one job
		err := pool.SubmitJob(job)
		if err != nil {
			panic(err)
		}
	}
	pool.Stop()
}
