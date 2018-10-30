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
	worker := NewWorker(pipeline, &Config{MaxRetries: 0, Concurrency: 2})
	worker.Start()
	defer worker.Close()

	expected := []int{}
	nb := 25

	for i := 0; i < nb; i++ {
		worker.AddJob(i)
		expected = append(expected, 4*i*i+1)
	}

	var results []int
	for i := 0; i < nb; i++ {
		l := <-worker.ReturnChan
		results = append(results, l.(int))
	}

	if !sameSlice(results, expected) {
		t.Errorf("Pipeline Parallel was incorrect, got: %v, want: %v.", results, expected)
	}

}

func sameSlice(x, y []int) bool {
	if len(x) != len(y) {
		return false
	}
	// create a map of string -> int
	diff := make(map[int]int, len(x))
	for _, _x := range x {
		// 0 value for int is 0, so just increment a counter for the string
		diff[_x]++

	}
	for _, _y := range y {
		// If the string _y is not in diff bail out early
		if _, ok := diff[_y]; !ok {

			return false
		}
		diff[_y]--
		if diff[_y] == 0 {
			delete(diff, _y)
		}

	}

	if len(diff) == 0 {
		return true
	}
	return false
}
