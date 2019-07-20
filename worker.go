package pool

import (
	"sync"
)

type Worker struct {
	id      int
	queue   <-chan Job
	results chan<- *Result
	done    <-chan bool
	wg      *sync.WaitGroup
}

func (w *Worker) Run() {
	for {
		select {
		case j := <-w.queue:
			if j == nil {
				continue
			}
			val, err := j.Run()
			r := &Result{
				Value: val,
				Error: err,
			}
			w.results <- r
			w.wg.Done()
		}
	}
}
