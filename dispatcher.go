package scheduling

import (
	"fmt"
	"os"
)

// CancelableJobQueue ...
type CancelableJobQueue chan CancelableJob

// Dispatcher ...
type Dispatcher struct {
	jobStore   JobStore
	jobQueue   JobQueue
	workerPool chan CancelableJobQueue
	workers    []Worker
	maxWorkers int
	quit       chan bool
}

// NewDispatcher ...
func NewDispatcher(jobStore JobStore, jobQueue JobQueue, maxWorkers int) *Dispatcher {
	return &Dispatcher{
		jobStore:   jobStore,
		jobQueue:   jobQueue,
		workerPool: make(chan CancelableJobQueue, maxWorkers),
		workers:    make([]Worker, maxWorkers),
		quit:       make(chan bool),
		maxWorkers: maxWorkers,
	}
}

// Start ...
func (d *Dispatcher) Start() {
	for i := 0; i < d.maxWorkers; i++ {
		worker := NewWorker(d.workerPool)
		worker.Start()
		d.workers[i] = worker
	}
	go d.dispatch()
}

// Stop ...
func (d *Dispatcher) Stop() {
	for _, worker := range d.workers {
		worker.Stop()
	}
	go func() {
		d.quit <- true
	}()
}

func (d *Dispatcher) dispatch() {
	for {
		select {
		case serializedJob := <-d.jobQueue:

			cancelableJob := NewCancelableJob(serializedJob)
			d.jobStore[cancelableJob.ID] = cancelableJob

			go func(job CancelableJob) {
				select {
				// find available worker
				case workerChannel := <-d.workerPool:
					// push job to worker
					workerChannel <- job

				case <-job.Context.Done():
					fmt.Fprintf(os.Stderr, "job canceled: %s", job.Context.Err())
				}
			}(cancelableJob)

		case <-d.quit:
			return
		}
	}
}
