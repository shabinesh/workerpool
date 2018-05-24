package workers

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
)

type Worker struct {
	pool     *WorkerPool
	ID       string
	ctx      context.Context
	jobQueue chan func() // incoming jobs
}

func (w Worker) stop(reason string) error {
	log.WithField("reason", reason).WithField("workerID", w.ID).Debug("channel is closed, worker quitting")
	close(w.jobQueue)
	return nil
}

func timeit(t time.Time) {
	log.WithField("timeelapsed", time.Now().Sub(t)).Info("completed")
}

func (w Worker) Run() error {
	log.WithFields(log.Fields{"workerId": w.ID}).Info("worker started")
	for {
		// register as avaliable
		if err := w.pool.addToPool(w.jobQueue); err != nil {
			return w.stop("orphaned")
		}

		select {
		case job, ok := <-w.jobQueue:
			if !ok {
				return w.stop("worker queue closed")
			}
			defer timeit(time.Now())
			job()
		case <-w.ctx.Done():
			return w.stop("signal")
		}
	}
}

func NewWorker(ctx context.Context, pool *WorkerPool, id string) *Worker {
	return &Worker{
		ID:       id,
		ctx:      ctx,
		pool:     pool,
		jobQueue: make(chan func()),
	}
}
