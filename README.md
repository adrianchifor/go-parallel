# go-parallel

Simple and customizable job queue/goroutine manager.

Create a pool of workers and execute jobs in the background and in parallel.

## Install

```
go get -u github.com/adrianchifor/go-parallel
```

## Usage

There are 3 predefined job pools you can choose from:

* SmallJobPool - 10 workers, 100 job queue size
* MediumJobPool - 50 workers, 500 job queue size
* LargeJobPool - 100 workers, 1000 job queue size

Or build your own:

``` go
jobPoolConfig := parallel.JobPoolConfig{
  WorkerCount:  30,
  JobQueueSize: 300,
}

jobPool := parallel.CustomJobPool(jobPoolConfig)
```

### Example

``` go
package main

import (
  "fmt"
  "time"

  "github.com/adrianchifor/go-parallel"
)

func main() {
  jobPool := parallel.SmallJobPool()
  defer jobPool.Close()

  for i := 1; i <= 50; i++ {
    i := i

    // If queue is full, AddJob will block until there's space
    jobPool.AddJob(func() {
      fmt.Print(i, " ")
      time.Sleep(500 * time.Millisecond)
    })
  }

  err := jobPool.Wait()
  if err != nil {
    fmt.Println("\nError:", err.Error())
  }

  fmt.Println("\nDone")
}
```

Output:

```
3 10 8 4 5 7 6 2 9 1 11 12 14 13 18 20 15 17 19 16 21 23 22 25 26 27 24 30 29 28 31 32 33 34 35 36 37 38 39 40 42 43 45 46 47 41 49 44 50 48
Done
```

### Example with timeout

``` go
package main

import (
  "context"
  "fmt"
  "time"

  "github.com/adrianchifor/go-parallel"
)

func main() {
  jobPool := parallel.SmallJobPool()
  defer jobPool.Close()

  for i := 1; i <= 50; i++ {
    i := i

    jobPool.AddJob(func() {
      fmt.Print(i, " ")
      time.Sleep(500 * time.Millisecond)
    })
  }

  ctx, cancel := context.WithTimeout(context.Background(), time.Second)
  defer cancel()

  err := jobPool.WaitContext(ctx)
  if err != nil {
    fmt.Println("\nError:", err.Error())
  }

  fmt.Println("\nDone")
}
```

Output:

```
7 8 5 9 6 3 2 1 10 4 11 12 13 14 15 16 17 19 18 20
Error: context deadline exceeded

Done
```
