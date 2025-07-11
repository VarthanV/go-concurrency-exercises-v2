# Incorrect synchronization

- Programs with races are incorrect and can exhibit non-sequentially consistent executions. In particular, note that a read r may observe the value written by any write w that executes concurrently with r. Even if this occurs, it does not imply that reads happening after r will observe writes that happened before w.

```go
var a, b int

func f() {
	a = 1
	b = 2
}

func g() {
	print(b)
	print(a)
}

func main() {
	go f()
	g()
}

```

- When it comes to programs with races, both programmers and compilers should remember the advice: don't be clever.


