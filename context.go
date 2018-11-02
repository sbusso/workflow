package workflow

import (
	"sync"
)

type Config struct {
	MaxRetries  int
	Concurrency int
}

type Context struct {
	wg        *sync.WaitGroup
	jobCount  int32
	doneCount int32
	workflow  *Workflow
}

func NewContext(workflow *Workflow) *Context {
	return &Context{
		wg:       new(sync.WaitGroup),
		workflow: workflow,
	}
}
