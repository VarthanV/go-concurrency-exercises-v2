package generators

import (
	"log"
	"time"
)

func SimpleGeneratorDriver() {
	done := make(chan any)
	defer close(done)
	s := Generate(done, 1, 2, 3, 4, 5)
	for v := range s {
		println(v)
	}
}

func SimpleGeneratorFromFuncDriver() {
	done := make(chan any)
	defer close(done)
	s := GenerateFromFunc(done, func() []string {
		time.Sleep(2 * time.Second) // Simulate some delay
		return []string{"a", "b", "c", "d", "e"}
	})
	log.Println("Generator started, waiting for values to be received...")
	for v := range s {
		println(v)
	}
}
