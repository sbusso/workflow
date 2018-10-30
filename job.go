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

func NewJob(ctx *Context, cfg *Config, data interface{}) *Job {
	ctx.wg.Add(1)
	ctx.jobCount++
	return &Job{Context: ctx, Config: cfg, Data: data}
}

func NewSerialJob(data interface{}) *Job {
	return &Job{Data: data}
}

func (j *Job) Error(err string) {
	j.Err = fmt.Errorf(err)
}

func (j *Job) Retry(idx int) {
	j.retries++
	if j.retries < j.MaxRetries {
		j.workflow.ReQueueJob(j, idx)
	} else {
		j.Done()
	}
}
