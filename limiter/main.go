package main

import (
	"context"
	"time"
)

func main() {
	l := NewNaive(3, time.Second)
	l.AddTasks([]Task{{URL: "https://jsonplaceholder.typicode.com/todos/1"}, {URL: "https://jsonplaceholder.typicode.com/todos/2"}})
	l.Start(context.Background())
}
