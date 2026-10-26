package worker

import (
	"context"
	"log"
	"sync"
)

type Job func(context.Context) error

type Pool struct {
	jobs   chan Job
	wg     sync.WaitGroup
	cancel context.CancelFunc
}

func NewPool(parent context.Context, workers, queueSize int) *Pool {
	if workers < 1 {
		workers = 1
	}
	if queueSize < 1 {
		queueSize = 1
	}

	ctx, cancel := context.WithCancel(parent)
	p := &Pool{
		jobs:   make(chan Job, queueSize),
		cancel: cancel,
	}

	p.wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			defer p.wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case job, ok := <-p.jobs:
					if !ok {
						return
					}
					if job != nil {
						if err := job(ctx); err != nil {
							log.Printf("worker job failed: %v", err)
						}
					}
				}
			}
		}()
	}

	return p
}

func (p *Pool) Submit(job Job) bool {
	if job == nil {
		return false
	}

	select {
	case p.jobs <- job:
		return true
	default:
		return false
	}
}

func (p *Pool) Close() {
	close(p.jobs)
	p.cancel()
	p.wg.Wait()
}
