package fanout

import (
	"log"
	"time"

	"github.com/VarthanV/golang-exercises-concurrency/generators"
)

func NaiveFanout() {
	done := make(chan any)
	defer close(done)
	in := generators.Generate(done, 1, 2, 3, 4, 5)
	out := Fanout(done, in, 3, func(val int) int {
		return val * 2 // Example work function that doubles the value
	})
	log.Println("Fanout started, waiting for values to be received...")
	for v := range out {
		println(v)
	}
}

func FanOutWithSem() {
	done := make(chan any)
	defer close(done)
	total := 100000
	var (
		integers = []int{}
	)

	for i := range total {
		integers = append(integers, i)
	}

	in := generators.Generate(done, integers...)
	out := FanoutWithSem(done, in, 100, 5, func(val int) int {
		time.Sleep(2 * time.Second)
		return val * 2 // Example work function that doubles the value
	})
	log.Println("Fanout with semaphore started, waiting for values to be received...")
	for v := range out {
		println(v)
	}
	log.Println("Fanout with semaphore completed.")
}
