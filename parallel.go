package parallel

import (
	"context"
	"errors"
	"sync"
)

type JobPool struct {
	jobsChan  chan func()
	waitGroup *sync.WaitGroup
}

type JobPoolConfig struct {
	WorkerCount  int
	JobQueueSize int
}

func CustomJobPool(config JobPoolConfig) *JobPool {
	jobPool := &JobPool{
		jobsChan:  make(chan func(), config.JobQueueSize),
		waitGroup: &sync.WaitGroup{},
	}

	for i := 1; i <= config.WorkerCount; i++ {
		go jobPool.runWorker()
	}

	return jobPool
}

func SmallJobPool() *JobPool {
	return CustomJobPool(JobPoolConfig{
		WorkerCount:  10,
		JobQueueSize: 100,
	})
}

func MediumJobPool() *JobPool {
	return CustomJobPool(JobPoolConfig{
		WorkerCount:  50,
		JobQueueSize: 500,
	})
}

func LargeJobPool() *JobPool {
	return CustomJobPool(JobPoolConfig{
		WorkerCount:  100,
		JobQueueSize: 1000,
	})
}

func (jp *JobPool) runWorker() {
	for work := range jp.jobsChan {
		work()
		jp.waitGroup.Done()
	}
}

func (jp *JobPool) AddJob(job func()) error {
	if job == nil {
		return errors.New("JobPool.AddJob: job function cannot be nil")
	}

	jp.jobsChan <- job
	jp.waitGroup.Add(1)

	return nil
}

func (jp *JobPool) WaitContext(ctx context.Context) error {
	channel := make(chan bool)
	go func() {
		jp.waitGroup.Wait()
		close(channel)
	}()

	select {
	case <-channel:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (jp *JobPool) Wait() error {
	return jp.WaitContext(context.Background())
}

func (jp *JobPool) Close() {
	close(jp.jobsChan)
}
