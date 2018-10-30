package workflow

import (
	"fmt"
)

type Worker struct {
	*Pipeline
	*Context
	*Config
}

func NewWorker(pipeline *Pipeline, cfg *Config) *Worker {

	// Config for Processors / Context for Worker Caller

	w := &Worker{
		Pipeline: pipeline,
		Config:   cfg,
	}

	ctx := NewContext(w, len(pipeline.steps))

	w.Context = ctx

	return w
}

func (w *Worker) AddJob(item interface{}) {
	go queue(NewJob(w.Context, item), w.chans[0])
}

func (w *Worker) ReQueueJob(job *Job) {
	go queue(job, w.chans[0])
}

// Start scraper workers
func (w *Worker) Start() {
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
func (w *Worker) Close() {
	for _, c := range w.chans {
		close(c)
	}
}

func (w *Worker) Next() {

}

// Queue a job in any channel queue
func queue(job *Job, ch chan *Job) {
	ch <- job
}

func (w *Worker) procWorker(proc Processor, idx int, next chan *Job) {
	for job := range w.chans[idx] {
		w.steps[idx](job)

		if job.Err == nil {
			go queue(job, next)
		} else {
			job.Error(fmt.Sprintf("Processor error %v\n", job.Err))
		}
	}
}

func (w *Worker) resultWorker() {
	for job := range w.resultChan {
		w.ReturnChan <- job.Data
		job.Done()
	}
}
