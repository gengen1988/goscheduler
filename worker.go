package scheduling

// Worker ...
type Worker struct {
	workerPool chan CancelableJobQueue
	jobChannel chan CancelableJob
	quit       chan bool
}

// NewWorker ...
func NewWorker(workerPool chan CancelableJobQueue) Worker {
	return Worker{
		workerPool: workerPool,
		jobChannel: make(chan CancelableJob),
		quit:       make(chan bool),
	}
}

// Start ...
func (w *Worker) Start() {
	go func() {
		for {
			w.workerPool <- w.jobChannel
			select {
			case cancelableJob := <-w.jobChannel:
				cancelableJob.Execute(cancelableJob.Context)

			case <-w.quit:
				return
			}
		}
	}()
}

// Stop ...
func (w Worker) Stop() {
	go func() {
		w.quit <- true
	}()
}
