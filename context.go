package workflow

import (
	"runtime"
	"sync"
)

type Config struct {
	MaxRetries  int
	Concurrency int
}

type Context struct {
	ReturnChan chan interface{}
	wg         *sync.WaitGroup
	jobCount   int
	doneCount  int
	workflow   *Workflow
	chans      []chan *Job
}

func NewContext(workflow *Workflow, nbProcs int) *Context {
	var chs = make([]chan *Job, nbProcs)
	for i := range chs {
		chs[i] = make(chan *Job, runtime.NumCPU()*workflow.Concurrency)
	}
	return &Context{
		ReturnChan: make(chan interface{}, runtime.NumCPU()*workflow.Concurrency),
		wg:         new(sync.WaitGroup),
		workflow:   workflow,
		chans:      chs,
	}
}
