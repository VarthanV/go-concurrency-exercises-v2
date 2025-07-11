package mutexes

import (
	"log"
	"sync"
)

var (
	iterations = 1000
)

func incrementer(wg *sync.WaitGroup, mu *sync.Mutex, count *int) {
	defer wg.Done()
	for range iterations {
		mu.Lock()
		log.Println("Incrementing..........")
		*count += 1
		mu.Unlock()
	}
}

func decrementer(wg *sync.WaitGroup, mu *sync.Mutex, count *int) {
	defer wg.Done()
	for range iterations - 10 {
		mu.Lock()
		log.Println("Decrementing..........")
		*count -= 1
		mu.Unlock()
	}
}

func SimpleMutexCounterDriver() {
	var (
		wg    sync.WaitGroup
		mu    sync.Mutex
		count = 0
	)

	wg.Add(2)
	log.Println("Starting simple mutex counter.....")
	go incrementer(&wg, &mu, &count)
	go decrementer(&wg, &mu, &count)
	wg.Wait()
	log.Println("Count is ", count)
}
