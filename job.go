package workflow

import (
	"fmt"
)

// Job is the entity to handle worker activity and data
type Job struct {
	retries int
	Err     error
	Data    interface{}
	Index   int // current processor index
	*Context
	*Config
}

func NewJob(ctx *Context, data interface{}) *Job {
	ctx.wg.Add(1)
	ctx.jobCount++
	return &Job{Context: ctx, Data: data}
}

func NewSerialJob(data interface{}) *Job {
	return &Job{Data: data}
}

func (j *Job) Error(err string) {
	j.retries++

	if j.retries < j.MaxRetries {
		j.worker.ReQueueJob(j)
	} else {
		j.Err = fmt.Errorf(err)
		j.Done()
	}
}
