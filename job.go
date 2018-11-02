package workflow

import "sync/atomic"

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
	atomic.AddInt32(&ctx.jobCount, 1)
	return &Job{Context: ctx, Config: cfg, Data: data}
}

func NewSerialJob(data interface{}) *Job {
	return &Job{Data: data}
}

func (j *Job) Retry(err error, idx int) {
	j.Err = err
	j.retries++
	if j.retries < j.MaxRetries {
		j.workflow.ReQueueJob(j, idx)
	} else {
		j.Done()
	}
}

func (j *Job) Done() {
	atomic.AddInt32(&j.doneCount, 1)
	j.wg.Done()
}
