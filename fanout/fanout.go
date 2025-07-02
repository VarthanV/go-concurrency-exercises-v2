package fanout

import (
	"log"
	"sync"
)

func worker[T any](
	wg *sync.WaitGroup,
	done chan any, id int, in <-chan T, out chan<- T,
	work func(T) T) {
	log.Println("Worker", id, "started")
	defer wg.Done()
	for val := range in {
		select {
		case <-done:
			return
		case out <- work(val):
			// Value sent to the output channel
		}
	}
}

func workerWithSem[T any](wg *sync.WaitGroup,
	done chan any, id int, in <-chan T, out chan<- T,
	work func(T) T, sem chan struct{}) {
	sem <- struct{}{} // Acquire a semaphore slot
	wg.Add(1)
	worker(wg, done, id, in, out, work)
	<-sem
}

func Fanout[T any](done chan any, in <-chan T, numWorkers int, work func(val T) T) <-chan T {
	out := make(chan T)
	var (
		wg sync.WaitGroup
	)
	for i := range numWorkers {
		wg.Add(1)
		go worker(&wg, done, i, in, out, work)
	}

	go func() {
		wg.Wait()
		// Cleanup once the workers are done
		close(out)
	}()
	return out
}

func FanoutWithSem[T any](done chan any,
	in <-chan T, numWorkers int,
	concurrencyLimit int, work func(val T) T) <-chan T {
	out := make(chan T)
	sem := make(chan struct{}, concurrencyLimit) // Semaphore to limit concurrency
	var (
		wg sync.WaitGroup
	)
	for i := range numWorkers {
		go workerWithSem(&wg, done, i, in, out, work, sem)
	}

	go func() {
		wg.Wait()
		// Cleanup once the workers are done
		close(out)
	}()
	return out
}
