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
	wg         *sync.WaitGroup
	jobCount   int
	doneCount  int
	workflow   *Workflow
	chans      []chan *Job
	ReturnChan chan interface{}
	resultChan chan *Job
}

func NewContext(workflow *Workflow, nbProcs int) *Context {
	var chs = make([]chan *Job, nbProcs)
	for i := range chs {
		chs[i] = make(chan *Job, runtime.NumCPU()*workflow.Concurrency)
	}
	return &Context{
		wg:         new(sync.WaitGroup),
		workflow:   workflow,
		chans:      chs,
		ReturnChan: make(chan interface{}, runtime.NumCPU()*workflow.Concurrency),
		resultChan: make(chan *Job, runtime.NumCPU()*workflow.Concurrency),
	}
}

func (ctx *Context) Done() {
	ctx.doneCount++
	ctx.wg.Done()
}
