package pool

import (
	"errors"
	"sync"
)

var (
	errAddWouldBlock = errors.New("Pool not started, Add would block")
	errPoolStarted   = errors.New("Pool already started")
	errPoolStopped   = errors.New("Pool has already stopped")
)

// A Pool is a concurrency-limited set of workers that can run any kind of Job.
type Pool struct {
	limit            int
	queue            chan Job
	resultValues     []*Result
	results          chan *Result
	started, stopped bool
	wg               *sync.WaitGroup
	workers          []*Worker
}

// New takes an Options and returns a new Pool.
func New(opts *Options) *Pool {
	return &Pool{
		limit:   opts.Limit,
		wg:      &sync.WaitGroup{},
		queue:   make(chan Job),
		results: make(chan *Result),
	}
}

// Start enables a pool to begin processing jobs. Returns an error if already
// started.
func (p *Pool) Start() error {
	if p.started {
		return errPoolStarted
	}
	p.started = true

	go p.collectResults()

	for i := 0; i < p.limit; i++ {
		w := &Worker{
			id:      i,
			queue:   p.queue,
			results: p.results,
			wg:      p.wg,
		}
		go w.Run()
		p.workers = append(p.workers, w)
	}
	return nil
}

// Add queues a Job for processing.
func (p *Pool) Add(j Job) error {
	if !p.started {
		return errAddWouldBlock
	}
	if p.stopped {
		return errPoolStopped
	}
	p.wg.Add(1)
	p.queue <- j
	return nil
}

func (p *Pool) collectResults() {
	for {
		select {
		case r := <-p.results:
			if r == nil {
				continue
			}
			p.resultValues = append(p.resultValues, r)
		}
	}
}

// GetResults waits for all jobs to finish processing and collects the results.
// Results are wrapped in a Result, and are not necessarily returned in the
// order they were queued.
func (p *Pool) GetResults() []*Result {
	p.stopped = true
	p.wg.Wait()
	defer close(p.queue)
	defer close(p.results)
	return p.resultValues
}
