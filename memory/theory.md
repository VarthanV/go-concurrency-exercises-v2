# Memory Model

## Introduction

The Go memory model specifies the conditions under which reads of a variable in one goroutine can be guaranteed to observe values produced by writes to the same variable in a different goroutine.


## Advice

- Programs that modify data being simultaneously accessed by multiple goroutines must serialize such access.

- To serialize access, protect the data with channel operations or other synchronization primitives such as those in the ``sync`` and ``sync/atomic`` packages.

## Informal Overview

- Go approaches its memory model in much the same way as the rest of the language, aiming to keep the semantics simple, understandable, and useful.

- A ``data race ``is defined as a write to a memory location happening concurrently with another read or write to that same location, unless all the accesses involved are atomic data accesses as provided by the sync/atomic package.

- In the absence of data races, Go programs behave as if all the goroutines were multiplexed onto a single processor.

## Memory Model

- The memory model describes the requirements on program executions, which are made up of goroutine executions, which in turn are made up of memory operations.

- A memory operation is modeled by four details:
    - its kind, indicating whether it is an ordinary data read, an ordinary data write, or a synchronizing operation such as an atomic data access, a mutex operation, or a channel operation,

    - its location in the program,

    - the memory location or variable being accessed, and

    - the values read or written by the operation.

- Some memory operations are read-like, including read, atomic read, mutex lock, and channel receive.

-  Other memory operations are write-like, including write, atomic write, mutex unlock, channel send, and channel close. Some, such as atomic compare-and-swap, are both read-like and write-like.

- A goroutine execution is modeled as a set of memory operations executed by a single goroutine.

**Requirement 1**:  The memory operations in each goroutine must correspond to a correct sequential execution of that goroutine, given the values read from and written to memory.

- That execution must be consistent with the sequenced before relation, defined as the partial order requirements set out by the Go language specification .


Sure!

In Go, **"sequenced before"** refers to the **guaranteed order of execution** between operations, as defined by the **Go language specification**.

### 🔹 In short:

> **If operation A is sequenced before B, then A is guaranteed to run before B.**

This ensures **predictable behavior** — especially important when reasoning about things like:

* Variable assignments
* Function calls
* Evaluation order of expressions

### ✅ Example (consistent with "sequenced before"):

```go
x := 1       // A
y := x + 1   // B
```

Here, A is sequenced before B — `x` is assigned before it's used.

---

If your program's execution **violates** the "sequenced before" relationship (e.g., due to data races or undefined evaluation order), the result is **unpredictable behavior**.

- A Go program execution is modeled as a set of goroutine executions, together with a mapping W that specifies the write-like operation that each read-like operation reads from. (Multiple executions of the same program can have different program executions.).

## Implementation restirctions for Programs containing Data Races

- Any implementation can, upon detecting a data race, report the race and halt execution of the program. Implementations using ThreadSanitizer (accessed with “go build -race”) do exactly this.

- A read of an array struct or complex number may be implemented as read of each individual subvalue(array element, struct field, or real/imaginary component) in any order. Similarly, a write of an array, struct, or complex number may be implemented as a write of each individual sub-value, in any order.

## Synchornization

- Program initialization runs in a single goroutine, but that goroutine may create other goroutines, which run concurrently.


- If package p imports q , then completion of q's init functions happens before start of any p.

- The completion of all ``init`` functions is synchronized before the start of the function ``main.main``.


## Goroutine creation

- The go statement that starts a new goroutine is synchornized before the start of goroutine's execution.

Eg

```go
var a string

func f() {
	print(a)
}

func hello() {
	a = "hello, world"
	go f()
}
```

calling hello will print "hello, world" at some point in the future (perhaps after hello has returned).

- It reflect the value "Hello world" since it is in the memory before kickstarting of the goroutine unless and until we have a mutex if we change value of a immediately after the go routine spawning it is not guaranteed to reflect.

| Action                               | Memory sync guaranteed?                     |
| ------------------------------------ | ------------------------------------------- |
| Writes **before** the `go` statement | ✅ Yes, visible to the new goroutine         |
| Writes **after** the `go` statement  | ❌ No, use sync/atomic, channels, or mutexes |


## Goroutine destruction

- The exit of a goroutine is not guaranteed to be synchronized before any event in the program. For example, in this program:

```go
var a string

func hello() {
	go func() { a = "hello" }()
	print(a)
}
```
- the assignment to a is not followed by any synchronization event, so it is not guaranteed to be observed by any other goroutine. In fact, an aggressive compiler might delete the entire go statement.

- If the effects of a goroutine must be observed by another goroutine, use a synchronization mechanism such as a lock or channel communication to establish a relative ordering.

## Channel communication

- Channel communication is the main method of synchronization between goroutines. Each send on a particular channel is matched to a corresponding receive from that channel, usually in a different goroutine.


- A send on a channel is synchronized before the completion of the corresponding receive from that channel.

```go
var c = make(chan int, 10)
var a string

func f() {
	a = "hello, world"
	c <- 0
}

func main() {
	go f()
	<-c
	print(a)
}
```
is guaranteed to print "hello, world". The write to a is sequenced before the send on c, which is synchronized before the corresponding receive on c completes, which is sequenced before the print.

- The closing of a channel is synchronized before a receive that returns a zero value because the channel is closed.

### ✅ Example: Shared variable synchronized by channel close

```go
package main

import (
	"fmt"
)

func main() {
	var shared int
	ch := make(chan struct{})
	done := make(chan struct{})

	// Writer goroutine: sets shared = 42, then closes channel
	go func() {
		shared = 42
		close(ch) // synchronization point
	}()

	// Reader goroutine: waits for channel to be closed, then reads shared
	go func() {
		<-ch // guarantees all memory writes before close() are visible
		fmt.Println("Shared value:", shared) // Guaranteed to print 42
		close(done)
	}()

	<-done
}
```

---

### 🔍 What’s Happening:

1. **Writer goroutine** sets `shared = 42`, **then closes** channel `ch`.
2. **Reader goroutine** waits (`<-ch`) for the close signal.
3. Go **guarantees**: all memory writes (like `shared = 42`) done before `close(ch)` are visible to the reader.

---

### 🧠 So Why Is This Safe?

This pattern works **without locks** because:

> `close(ch)` in one goroutine and `<-ch` in another **synchronize memory** — all writes before `close` are visible after `<-ch`.

If this channel synchronization were not there, you might get a **data race** and see the wrong value (like 0 instead of 42).


- A receive from an unbuffered channel is synchronized before the completion of the corresponding send on that channel.

```go
var c = make(chan int)
var a string

func f() {
	a = "hello, world"
	<-c
}

func main() {
	go f()
	c <- 0
	print(a)
}
```

- is also guaranteed to print "hello, world". The write to a is sequenced before the receive on c, which is synchronized before the corresponding send on c completes, which is sequenced before the print.

- If the channel were ``buffered`` ``(e.g., c = make(chan int, 1))`` then the program ``would not be guaranteed to print "hello, world"``. (It might print the empty string, crash, or do something else.)

- The ``kth`` receive from a channel with capacity ``C`` is synchronized before the completion of the k+Cth send on that channel.


- In other words 
    - The channel buffer only has ``C`` slots.
    - Once the buffer is full the next send (K+cth) send must wait for an earlier receive kth receive to finish and free up space.
    - That means the kth receive happens-before the k+Cth send completes.

**Visualization**
- Let’s say the channel has ``C = 2``.
- We do 3 sends
- First send → goes into buffer slot 1.
- Second send -> goes to buffer slot 2.
- Third send → now the buffer is full, so this blocks until the first receive happens.
- So the 1st receive (k=1) must complete before the 3rd send (k+C=3) can proceed.



- This is a synchronization point in Go — even though the goroutines don’t explicitly use locks, the channel buffer implicitly enforces ordering. 


## Locks

- The sync package implements two lock data types, ``sync.Mutex ``and ``sync.RWMutex``.

- For any sync.Mutex or sync.RWMutex variable l and n < m, call n of l.Unlock() is synchronized before call m of l.Lock() returns.

```go
var l sync.Mutex
var a string

func f() {
	a = "hello, world"
	l.Unlock()
}

func main() {
	l.Lock()
	go f()
	l.Lock()
	print(a)
}
```

- is guaranteed to print "hello, world". The first call to l.Unlock() (in f) is synchronized before the second call to l.Lock() (in main) returns, which is sequenced before the print.

- A successful call to ``l.TryLock (or l.TryRLock)`` is equivalent to a call to`` l.Lock`` (or l.RLock). An unsuccessful call has no synchronizing effect at all. As far as the memory model is concerned, l.TryLock (or l.TryRLock) may be considered to be able to return false even when the mutex l is unlocked


- If one goroutine unlocks a sync.Mutex, and another goroutine locks it afterward. Everything done before ``Unlock()`` **is guaranteed to be visible** after ``Lock()`` returns.

```go
var mu sync.Mutex
var msg string

func f() {
    msg = "hello, world" // write shared data
    mu.Unlock()          // release lock
}

func main() {
    mu.Lock()            // take lock (will block until unlocked in f)
    go f()               // start f
    mu.Lock()            // wait here until f unlocks
    fmt.Println(msg)     // guaranteed to see "hello, world"
}

```
- A read lock (``RLock``) will wait for any active write lock (Lo``ck) to be unlocked (Unlock).

- The write lock's Unlock() happens before the reader's RLock() returns.

- Similarly, once the reader is done and calls RUnlock(), then the next writer’s Lock() will only return after that reader finishes.

- Writers block readers

- Readers block writers (but not other readers).

- The memory model guarantees that this order is safe and consistent.


| Action               | Synchronization Effect | Safe for Shared Memory?           |
| -------------------- | ---------------------- | --------------------------------- |
| `Unlock → Lock`      | Yes                    | ✅                                 |
| `Unlock → RLock`     | Yes                    | ✅                                 |
| `RUnlock → Lock`     | Yes                    | ✅                                 |
| Successful `TryLock` | Yes (like `Lock`)      | ✅                                 |
| Failed `TryLock`     | No sync                | ❌ Not safe for memory assumptions |


