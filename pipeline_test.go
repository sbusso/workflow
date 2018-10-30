package workflow

import (
	"testing"
)

func double(j *Job) {
	v := j.Data.(int)
	j.Data = v * 2
}

func sq(j *Job) {
	v := j.Data.(int)
	j.Data = v * v
}

func up(j *Job) {
	v := j.Data.(int)
	j.Data = v + 1
}

func TestPipeline(t *testing.T) {
	pipeline := NewPipeline(double, sq, up)
	if l := len(pipeline.steps); l != 3 {
		t.Errorf("Pipeline Len was incorrect, got: %d, want: %d.", l, 3)
	}
}

func TestPipelineSerial(t *testing.T) {
	pipeline := NewPipeline(double, sq, up)

	if l := pipeline.Exec(2); l != 17 {
		t.Errorf("Pipeline Serial was incorrect, got: %d, want: %d.", l, 17)
	}

}

func TestPipelineParallel(t *testing.T) {
	pipeline := NewPipeline(double, sq, up)
	returnChan := make(chan interface{}, 1)
	worker := NewWorker(pipeline, returnChan, &Config{maxRetries: 0, concurrency: 2})
	worker.Start()
	worker.AddJob(2)
	defer worker.Close()

	if l := <-returnChan; l != 17 {
		t.Errorf("Pipeline Parallel was incorrect, got: %d, want: %d.", l, 17)
	}

}
