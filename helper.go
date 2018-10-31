package workflow

func ChanResultHelper(returnChan chan interface{}) Processor {
	return func(job *Job) {
		returnChan <- job.Data
	}
}
