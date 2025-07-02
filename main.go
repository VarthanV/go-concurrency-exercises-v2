package main

import (
	"github.com/VarthanV/golang-exercises-concurrency/fanout"
	"github.com/VarthanV/golang-exercises-concurrency/generators"
)

func main() {
	generators.SimpleGeneratorDriver()
	generators.SimpleGeneratorFromFuncDriver()

	fanout.NaiveFanout()
	fanout.FanOutWithSem()

}
