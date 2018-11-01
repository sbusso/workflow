package workflow

import (
	"runtime"
)

// Processor is a type of function to process document extracted from fetcher and produce output
type Processor func(interface{}) (interface{}, error)

type Workflow struct {
	steps []Processor
	chans []chan *Job

	*Config
	*Context
}

// SerialWorkflow
// ParallelWorkflow
// PersistentWorkflow

func NewWorkflow(cfg *Config, fs ...Processor) *Workflow {

	w := &Workflow{
		steps:  []Processor{},
		chans:  []chan *Job{},
		Config: cfg,
	}

	ctx := NewContext(w)
	w.AddProcessors(fs...)
	w.Context = ctx

	return w
}

func (w *Workflow) AddProcessors(fs ...Processor) {
	var chs = make([]chan *Job, len(fs))
	for i := range chs {
		chs[i] = make(chan *Job, runtime.NumCPU()*w.Concurrency)
	}
	w.steps = append(w.steps, fs...)
	w.chans = append(w.chans, chs...)
}

func (w *Workflow) Exec(data interface{}) interface{} {
	for _, proc := range w.steps {
		data, _ = proc(data)
	}
	return data
}

func (w *Workflow) AddJob(data interface{}) {
	go queue(NewJob(w.Context, w.Config, data), w.chans[0])
}

func (w *Workflow) ReQueueJob(job *Job, idx int) {
	go queue(job, w.chans[idx])
}

// Start scraper workers
func (w *Workflow) Start() {
	l := len(w.steps)
	for wk := 0; wk < w.Concurrency; wk++ {
		for i, p := range w.steps {
			var next chan *Job
			if i < l-1 {
				next = w.chans[i+1]
			} else {
				next = nil
			}
			go w.procWorker(p, i, next)
		}
	}
}

func (w *Workflow) Wait() {
	w.wg.Wait()
}

// Close stop all channels and shutdown workers
func (w *Workflow) Close() {
	for _, c := range w.chans {
		close(c)
	}
}

// func (w *Workflow) Next() {
// }

// Queue a job in any channel queue
func queue(job *Job, ch chan *Job) {
	ch <- job
}

func (w *Workflow) procWorker(proc Processor, idx int, next chan *Job) {
	for job := range w.chans[idx] {
		var err error
		d, err := w.steps[idx](job.Data)

		if err == nil {
			job.Data = d
			if next != nil {
				go queue(job, next)
			} else {
				go job.Done()
			}

		} else {
			go job.Retry(err, idx)
		}
	}
}
