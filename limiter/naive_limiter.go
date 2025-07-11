package main

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"
)

type Action func(a ...any) any

type Task struct {
	URL string
}

type naiveLimitedRunner struct {
	// semaphore to grant access
	sem chan struct{}
	// No of ops that  can happen concurrently in the given interval
	concurrencyLimit int
	// The interval at which the limit can reset
	resetInterval time.Duration
	// Assuming we only fetching of urls only can modify later when needed
	tasks []Task
	// waitgroup counter to implement the fork join cycle
	wg *sync.WaitGroup
}

func NewNaive(limit int, resetInterval time.Duration) *naiveLimitedRunner {
	l := &naiveLimitedRunner{
		concurrencyLimit: limit,
		sem:              make(chan struct{}, limit),
		resetInterval:    resetInterval,
		wg:               &sync.WaitGroup{},
	}

	for range limit {
		l.sem <- struct{}{}
	}
	return l
}

func (n *naiveLimitedRunner) AddTasks(t []Task) {
	n.tasks = t
}

func (n *naiveLimitedRunner) taskGenerator(ctx context.Context) <-chan Task {
	t := make(chan Task)
	go func() {
		// When no tasks is there the stream closes
		defer close(t)
		for _, task := range n.tasks {
			select {
			case <-ctx.Done():
				return
			case t <- task:
			}
		}
	}()

	return t
}

func (n *naiveLimitedRunner) Start(ctx context.Context) error {
	var (
		// We can spawn upto maxWorkers , but only  concurrencyLimit workers
		// must perform an action
		maxWorkers = 10000
	)

	taskStream := n.taskGenerator(ctx)

	go n.resetLimit(ctx)
	log.Println("Spawning workers ", maxWorkers)
	log.Println("Allowed concurrency limit ", n.concurrencyLimit)

	for range maxWorkers {
		n.wg.Add(1)
		go n.doAction(ctx, taskStream)
	}

	n.wg.Wait()
	close(n.sem)

	return nil
}

func (n *naiveLimitedRunner) resetLimit(ctx context.Context) {
	ticker := time.NewTicker(n.resetInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case _, ok := <-n.sem:
			if !ok {
				return
			}
		case <-ticker.C:
			log.Println("resetting the concurrency limit to ", n.concurrencyLimit)
			// Drain and fill
			for {
				select {
				case <-ctx.Done():
					return

				case <-ticker.C:
					log.Println("resetting the concurrency limit to", n.concurrencyLimit)

					// Drain
					for {
						select {
						case <-n.sem:
							// continue draining
						default:
							// channel is empty
							goto refill
						}
					}

				refill:
					for i := 0; i < n.concurrencyLimit; i++ {
						n.sem <- struct{}{}
					}
				}
			}
		}
	}
}

func (n *naiveLimitedRunner) doAction(ctx context.Context, stream <-chan Task) {
	defer n.wg.Done()
	for t := range stream {
		select {
		case <-ctx.Done():
			return
		default:
			// Acquire the token
			log.Println("got task ", t)
			<-n.sem
			n.fetch(ctx, t.URL)
			// Releasing is optional here , because the token  filling must
			// be handled wrt to the reset interval and concurrency limit
		}

	}
}

func (n *naiveLimitedRunner) fetch(ctx context.Context, url string) {
	log.Println("Fetching url ", url)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Println("error in creating req ", err)
		return
	}
	log.Println("Doing req")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("error in doing request ", err)
		return
	}
	defer res.Body.Close()

	log.Println("Status code ", res.StatusCode)
}
