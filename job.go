package scheduling

import "context"

// Job ...
type Job interface {
	Execute(context.Context)
}

// SerializedJob ...
type SerializedJob struct {
	ID string
	Job
}

// NewSerializedJob ...
func NewSerializedJob(id string, job Job) SerializedJob {
	return SerializedJob{id, job}
}

// CancelableJob ...
type CancelableJob struct {
	SerializedJob
	Cancel  context.CancelFunc
	Context context.Context
}

// NewCancelableJob ...
func NewCancelableJob(origin SerializedJob) CancelableJob {
	ctx, cancel := context.WithCancel(context.Background())
	return CancelableJob{origin, cancel, ctx}
}

// JobStore ...
type JobStore map[string]CancelableJob

// JobQueue ...
type JobQueue chan SerializedJob
