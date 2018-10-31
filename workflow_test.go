package workflow

import (
	"testing"
	"time"
)

func double(data interface{}) (interface{}, error) {
	d := data.(int)
	return d * 2, nil
}

func sq(data interface{}) (interface{}, error) {
	d := data.(int)
	return d * d, nil
}

func up(data interface{}) (interface{}, error) {
	d := data.(int)
	return d + 1, nil
}

func buildExpectation(n int) []int {
	expected := []int{}
	for i := 0; i < n; i++ {
		expected = append(expected, 4*i*i+1)
	}
	return expected
}
func TestWorkflow(t *testing.T) {
	workflow := NewWorkflow(&Config{}, double, sq, up)
	if l := len(workflow.steps); l != 3 {
		t.Errorf("Workflow Len was incorrect, got: %d, want: %d.", l, 3)
	}
}

func TestWorkflowSerial(t *testing.T) {
	workflow := NewWorkflow(&Config{}, double, sq, up)
	expected := 17
	if got := workflow.Exec(2).(int); got != expected {
		t.Errorf("Workflow Serial was incorrect, got: %d, want: %d.", got, expected)
	}
}

func TestWorkflowParallel(t *testing.T) {
	nb := 25

	workflow := NewWorkflow(&Config{MaxRetries: 0, Concurrency: 2}, double, sq, up)
	returnChan := workflow.HelperAddChanResult(nb)
	workflow.Start()
	defer workflow.Close()

	expected := buildExpectation(nb)

	for i := 0; i < nb; i++ {
		workflow.AddJob(i)
	}

	var results []int
	for i := 0; i < nb; i++ {
		l := <-returnChan
		results = append(results, l.(int))
	}

	if !sameSlice(results, expected) {
		t.Errorf("Workflow Parallel was incorrect, got: %v, want: %v.", results, expected)
	}

}

func TestWait(t *testing.T) {
	nb := 20
	workflow := NewWorkflow(&Config{MaxRetries: 0, Concurrency: 2}, double, sq, up)
	returnChan := workflow.HelperAddChanResult(nb)
	workflow.Start()
	waitCh := make(chan struct{})
	defer workflow.Close()

	expected := []int{}

	for i := 0; i < nb; i++ {
		workflow.AddJob(i)
		expected = append(expected, 4*i*i+1)
	}

	go func() {
		for {
			<-returnChan
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

func step1(w *Workflow) Processor {
	return func(data interface{}) (interface{}, error) {
		v := data.(int)

		if v < 5 {
			w.AddJob(v + 1)
		}

		return v + 1, nil
	}
}

func step2(data interface{}) (interface{}, error) {
	d := data.(int)
	return d * d, nil
}

func TestSpawning(t *testing.T) {
	workflow := NewWorkflow(&Config{MaxRetries: 0, Concurrency: 2})
	workflow.AddProcessors(step1(workflow), step2)
	returnChan := workflow.HelperAddChanResult(5)
	workflow.Start()
	waitCh := make(chan struct{})
	defer func() {
		workflow.Close()
		close(returnChan)
	}()

	expected := []int{4, 9, 16, 25, 36}
	got := []int{}

	workflow.AddJob(1)

	go func() {
		// for mdg := range returnChan
		for {
			msg, open := <-returnChan
			if !open {
				break
			}
			got = append(got, msg.(int))
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
	<-time.After(100 * time.Millisecond)
	if !sameSlice(got, expected) {
		t.Errorf("TestSpawning was incorrect, got: %v, want: %v.", got, expected)
	}
}

func sameSlice(x, y []int) bool {
	if len(x) != len(y) {
		return false
	}

	diff := make(map[int]int, len(x))
	for _, _x := range x {
		diff[_x]++
	}
	for _, _y := range y {
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
