package worker

// Work is Worker^s work interface
type Work interface {
	Execute()
}

// Worker is worker struct
type Worker struct {
	works []Work
}

// NewWorker create Worker
func NewWorker() *Worker {
	return &Worker{
		works: make([]Work, 0),
	}
}

// AddWork is add Work to Worker
func (w *Worker) AddWork(work Work) {
	w.works = append(w.works, work)
}

// Work call Execute method in Worker.works
func (w *Worker) Work() {
	for _, work := range w.works {
		work.Execute()
	}
}
