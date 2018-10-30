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
	worker     *Worker
	chans      []chan *Job
	ReturnChan chan interface{}
	resultChan chan *Job
}

func NewContext(worker *Worker, nbProcs int) *Context {
	var chs = make([]chan *Job, nbProcs)
	for i := range chs {
		chs[i] = make(chan *Job, runtime.NumCPU()*worker.Concurrency)
	}
	return &Context{
		wg:         new(sync.WaitGroup),
		worker:     worker,
		chans:      chs,
		ReturnChan: make(chan interface{}, runtime.NumCPU()*worker.Concurrency),
		resultChan: make(chan *Job, runtime.NumCPU()*worker.Concurrency),
	}
}

func (ctx *Context) Done() {
	ctx.doneCount++
	ctx.wg.Done()
}
