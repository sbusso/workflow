package workflow

import (
	"testing"
	"time"
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

func TestWorkflow(t *testing.T) {
	workflow := NewWorkflow(&Config{}, double, sq, up)
	if l := len(workflow.steps); l != 3 {
		t.Errorf("Workflow Len was incorrect, got: %d, want: %d.", l, 3)
	}
}

func TestWorkflowSerial(t *testing.T) {
	workflow := NewWorkflow(&Config{}, double, sq, up)

	if l := workflow.Exec(2); l != 17 {
		t.Errorf("Workflow Serial was incorrect, got: %d, want: %d.", l, 17)
	}
}

func TestWorkflowParallel(t *testing.T) {
	workflow := NewWorkflow(&Config{MaxRetries: 0, Concurrency: 2}, double, sq, up)
	workflow.Start()
	defer workflow.Close()

	expected := []int{}
	nb := 25

	for i := 0; i < nb; i++ {
		workflow.AddJob(i)
		expected = append(expected, 4*i*i+1)
	}

	var results []int
	for i := 0; i < nb; i++ {
		l := <-workflow.ReturnChan
		results = append(results, l.(int))
	}

	if !sameSlice(results, expected) {
		t.Errorf("Workflow Parallel was incorrect, got: %v, want: %v.", results, expected)
	}

}

func TestWait(t *testing.T) {
	workflow := NewWorkflow(&Config{MaxRetries: 0, Concurrency: 2}, double, sq, up)
	workflow.Start()
	waitCh := make(chan struct{})
	defer workflow.Close()

	expected := []int{}
	nb := 20

	for i := 0; i < nb; i++ {
		workflow.AddJob(i)
		expected = append(expected, 4*i*i+1)
	}

	go func() {
		for i := 0; i < nb; i++ {
			<-workflow.ReturnChan

		}
	}()

	go func() {
		workflow.Wait()
		close(waitCh)
	}()

	select {
	case <-waitCh:

	case <-time.After(100 * time.Millisecond):
		t.Error("Workflow Parallel didnt terminated by WaitGroup but with Timeout.")
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
