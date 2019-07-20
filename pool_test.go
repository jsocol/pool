package pool

import (
	"errors"
	"testing"
	"time"
)

type testJob struct {
	s string
	e error
}

func (j *testJob) Run() (interface{}, error) {
	time.Sleep(50 * time.Millisecond)
	if j.e != nil {
		return nil, j.e
	}
	return j.s, nil
}

func TestPool(t *testing.T) {
	p := New(&Options{
		Limit: 4,
	})

	jobs := []*testJob{
		{
			s: "hello",
		},
		{
			s: "goodbye",
		},
		{
			e: errors.New("whoops"),
		},
	}

	p.Start()

	for _, j := range jobs {
		p.Add(j)
	}

	results := p.GetResults()
	if len(results) != len(jobs) {
		t.Errorf("Wrong number of results, expected %d, got %d", len(jobs), len(results))
	}
}
