package workflow

import "fmt"

// Processor is a type of function to process document extracted from fetcher and produce output
type Processor func(*Job)

type Workflow struct {
	steps []Processor
	*Config
	*Context
}

// SerialWorkflow
// ParallelWorkflow
// PersistentWorkflow

func NewWorkflow(cfg *Config, fs ...Processor) *Workflow {

	w := &Workflow{
		steps:  fs,
		Config: cfg,
	}

	ctx := NewContext(w, len(fs))

	w.Context = ctx

	return w
}

func (w *Workflow) Exec(data interface{}) interface{} {
	j := NewSerialJob(data)
	for _, s := range w.steps {
		s(j)
	}
	return j.Data
}

func (w *Workflow) AddJob(item interface{}) {
	go queue(NewJob(w.Context, item), w.chans[0])
}

func (w *Workflow) ReQueueJob(job *Job) {
	go queue(job, w.chans[0])
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
				next = w.resultChan
			}
			go w.procWorker(p, i, next)
		}
	}

	go w.resultWorker()
}

// Close stop all channels and shutdown workers
func (w *Workflow) Close() {
	for _, c := range w.chans {
		close(c)
	}
}

func (w *Workflow) Next() {

}

// Queue a job in any channel queue
func queue(job *Job, ch chan *Job) {
	ch <- job
}

func (w *Workflow) procWorker(proc Processor, idx int, next chan *Job) {
	for job := range w.chans[idx] {
		w.steps[idx](job)

		if job.Err == nil {
			go queue(job, next)
		} else {
			job.Error(fmt.Sprintf("Processor error %v\n", job.Err))
		}
	}
}

func (w *Workflow) resultWorker() {
	for job := range w.resultChan {
		w.ReturnChan <- job.Data
		job.Done()
	}
}
