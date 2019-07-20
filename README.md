# Pool

**This is not production-ready software. Use at your own risk!**

I needed a quick way to be able to run a few jobs in parallel and
collect the results, with a limit on the concurrency. I found lots of
examples but not a lot of libraries, so I threw this together and put it
into a module.

It is not ready for production use. It is not well-tested, there are
probably issues.

Each pool can handle an arbitrary type of job, but it should probably be
thrown away after any given step, since it will be stopped when
`GetResults` is called.

## Basic Usage

```go
package main

import (
    "fmt"
    "time"

    "github.com/jsocol/pool"
)

// The structure of jobs is entirely up to you, however...
type job struct {
    delay int
    id    string
}

// ...jobs must implement this method signature
func (j *job) Run() (interface{}, error) {
    time.sleep(j.delay * time.Millisecond)
    return j.id, nil
}

func main() {
    // To create a new Pool, pass in an Options pointer. Right now Limit
    // is the only option
    p := pool.New(&pool.Options{
        Limit: 4,
    })

    // The Pool must be started before jobs can be added to it
    p.Start()

    for i := 0; i < 10; i++ {
        // Add jobs into the queue, they'll start processing right away
        p.Add(job{
            delay: i * 10,
            id:    fmt.Sprintf("job %d", i),
        })
    }

    // Wait for all jobs to process and get a []*pool.Result back.
    results := p.GetResults()
    for _, r := range results {
        if r.Error != nil {
            fmt.Println(r.Error)
        }
        // Convert r.Value back to whatever type your Run() method
        // returned
        fmt.Println(r.Value.(string))
    }
}
```
