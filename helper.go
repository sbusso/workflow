package workflow

func (w *Workflow) HelperAddChanResult(nb int) chan interface{} {
	returnChan := make(chan interface{}, nb)
	proc := func(data interface{}) (interface{}, error) {
		returnChan <- data
		return nil, nil
	}
	w.AddProcessors(proc)
	return returnChan
}
