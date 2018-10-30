package workflow

import (
	"runtime"
	"sync"
)

type Config struct {
	maxRetries  int
	concurrency int
}

type Context struct {
	wg         *sync.WaitGroup
	jobCount   int
	doneCount  int
	worker     *Worker
	chans      []chan *Job
	returnChan chan interface{}
	resultChan chan *Job
}

func NewContext(worker *Worker, nbProcs int, returnChan chan interface{}) *Context {
	var chs = make([]chan *Job, nbProcs)
	for i := range chs {
		chs[i] = make(chan *Job, runtime.NumCPU()*worker.concurrency)
	}
	return &Context{
		wg:         new(sync.WaitGroup),
		worker:     worker,
		chans:      chs,
		returnChan: returnChan,
		resultChan: make(chan *Job, runtime.NumCPU()*worker.concurrency),
	}
}

func (ctx *Context) Done() {
	ctx.doneCount++
	ctx.wg.Done()
}
