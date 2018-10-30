package workflow

// Processor is a type of function to process document extracted from fetcher and produce output
type Processor func(*Job)

type Pipeline struct {
	steps []Processor
}

func NewPipeline(fs ...Processor) *Pipeline {
	return &Pipeline{
		steps: fs,
	}
}

func (p *Pipeline) Exec(data interface{}) interface{} {
	j := NewSerialJob(data)
	for _, s := range p.steps {
		s(j)
	}
	return j.Data
}
