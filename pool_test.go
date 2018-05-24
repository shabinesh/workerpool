package workers

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stvp/assert"
)

func TestWorkerStarted(t *testing.T) {
	wp := NewWorkerPool(context.TODO(), 5)
	assert.Equal(t, wp.GetAvailability(), 5)
	wp.Stop()
	assert.Equal(t, wp.GetAvailability(), 0)
}

func TestPoolSendJobs(t *testing.T) {
	job := func() {
		fmt.Print("dsf")
	}
	wp := NewWorkerPool(context.TODO(), 1)
	defer wp.Stop()
	time.Sleep(2 * time.Second)
	wp.SubmitJob(job)
	wp.Stop()
}
