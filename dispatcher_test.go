package scheduling

import (
	"context"
	"fmt"
	"testing"
	"time"
)

type TestJob struct {
}

func (j *TestJob) Execute(ctx context.Context) {
	fmt.Println("execute job")

	select {
	case <-time.After(2 * time.Second):
		fmt.Println("job done")

	case <-ctx.Done():
		fmt.Println("job canceled")
	}
}

func NewTestJob() Job {
	return &TestJob{}
}

func TestMain(t *testing.T) {
	jobQueue := make(chan SerializedJob)
	jobStore := make(JobStore)
	dispatcher := NewDispatcher(jobStore, jobQueue, 1)

	dispatcher.Start()

	go func() {
		job1 := NewSerializedJob("123", &TestJob{})
		jobQueue <- job1
		job2 := NewSerializedJob("456", &TestJob{})
		jobQueue <- job2
	}()

	<-time.After(1 * time.Second)

	job, ok := jobStore["456"]
	if !ok {
		panic("can't find job")
	}
	job.Cancel()

	<-time.After(3 * time.Second)
	dispatcher.Stop()
}
