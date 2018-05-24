package workers

import (
	"context"
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
}

var (
	// ErrQueueFull when no avaliable worker are present
	ErrQueueFull = fmt.Errorf("queue full, try again")
)

// WorkerPool manages the workers
type WorkerPool struct {
	n                int
	jobQueue         chan func()
	avaliableWorkers chan chan func()
	mu               sync.Mutex
	ctx              context.Context
	cancel           context.CancelFunc
	stopped          bool
}

func (wp *WorkerPool) addToPool(task chan func()) error {
	if wp.stopped {
		return fmt.Errorf("stopped")
	}
	logrus.Debug("adding to pool")
	wp.avaliableWorkers <- task
	return nil
}

func (wp *WorkerPool) GetAvailability() int {
	if wp.stopped {
		return 0
	}

	return len(wp.avaliableWorkers)
}

func (wp *WorkerPool) Stop() {
	wp.stopped = true
	wp.cancel()
}

func (wp *WorkerPool) SubmitJob(job func()) error {
	wp.mu.Lock()
	defer wp.mu.Unlock()
	if wp.GetAvailability() > 0 {
		w := <-wp.avaliableWorkers
		w <- job
		logrus.Info("job submitted to worker %s", w)
		return nil
	}
	return ErrQueueFull
}

// NewWorkerPool creates a job queue and starts workers
func NewWorkerPool(ctx context.Context, n int) *WorkerPool {
	ctx, cancel := context.WithCancel(ctx)
	wp := &WorkerPool{
		n:                n,
		jobQueue:         make(chan func()),
		avaliableWorkers: make(chan chan func(), 100),
		cancel:           cancel,
		ctx:              ctx,
		stopped:          false,
	}

	for i := 0; i < wp.n; i++ {
		wid := fmt.Sprintf("worker-%d", i)
		worker := NewWorker(ctx, wp, wid)

		go worker.Run()

		if err := wp.addToPool(worker.jobQueue); err != nil {
			logrus.WithError(err).Debug("error while adding to avalibility queue")
			continue
		}
		fmt.Printf("worker %s started\n", wid)
	}

	return wp
}
