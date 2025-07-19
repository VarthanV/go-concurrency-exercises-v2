# Concurrency

## Share by communicating

- Concurrent programming in many envs is made difficult by subtelities required to implement correct access across shared variables.

- Go uses different approach of passing shared variables around channels and never actively shared by seperate threads of execution.

- Only one goroutine has access to the value at any given time, Data race cannot occur by design.

> Do not communicate by sharing memory; instead, share memory by communicating.

- Consider two independent single threaded programs, Now let these two communicate , if communication is synchornizer, there is still no need for synchornization.

- Unix pipelines fit the model perfectly, Although Go's approach to concurrency originates in Hoare's Communicating Sequential Processes (CSP), it can also be seen as a type-safe generalization of Unix pipes.

## Goroutines

- Goroutines have a simple model , it is a function executing concurrently as other goroutines in the same address space.

- It is a lightweight,costing little more than the allocation of a stack space. And the stacks are small so they start cheap ,and grow by allocating (and freeing) heap storage as required.

- Goroutines are multiplexed into multiple OS threads, If one should block because of waiting for I/O others continue to run. 

- To call a new goroutine prefix a function call with ``go``

- When the call completes, the goroutine exits, silently. (The effect is similar to the Unix shell's & notation for running a command in the background.)

```go
go list.Sort()  // run list.Sort concurrently; don't wait for it.
```

- A function literal can be handy in a goroutine invocation.

```go
func Announce(message string, delay time.Duration) {
    go func() {
        time.Sleep(delay)
        fmt.Println(message)
    }()  // Note the parentheses - must call the function.
}
```

- In Go, function literals are closures: the implementation makes sure the variables referred to by the function survive as long as they are active.

## Channels

- Like maps channels are allocated with make and resulting value acts as reference to an underlying data structure. 

- If an optional integer value is passed , it sets the capacity as buffer size of the channel.

- The default is 0 for synchornous or unbuffered channel

```go
ci := make(chan int)            // unbuffered channel of integers
cj := make(chan int, 0)         // unbuffered channel of integers
cs := make(chan *os.File, 100)  // buffered channel of pointers to Files
```

- Unbuffered channels combine communication—the exchange of a value—with synchronization—guaranteeing that two calculations (goroutines) are in a known state.

- A channel can allow the launching goroutine to wait for the sort to complete.

```go
    c := make(chan int)  // Allocate a channel.
    // Start the sort in a goroutine; when it completes, signal on the channel.
    go func() {
        list.Sort()
        c <- 1  // Send a signal; value does not matter.
    }()
    doSomethingForAWhile()
    <-c   // Wait for sort to finish; discard sent value.
```

- Receivers always block until there is data to receive. 

- If the channel is unbuffered, the sender blocks until the receiver receives data.

- If a channel is buffered the sender blocks only until the value has been copied to buffer; if the buffer is full this means the sender will be waiting until one of the block finishes.

- A buffered channel can be used as a semaphore to limit throughput, 

```go
var sem = make(chan int, MaxOutstanding)

func handle(r *Request) {
    sem <- 1    // Wait for active queue to drain.
    process(r)  // May take a long time.
    <-sem       // Done; enable next request to run.
}

func Serve(queue chan *Request) {
    for {
        req := <-queue
        go handle(req)  // Don't wait for handle to finish.
    }
}
```

- The capacity of the channel limits the number of calls to process.

- This design has a problem, though: Serve creates a new goroutine for every incoming request, even though only MaxOutstanding of them can run at any moment. As a result, the program can consume unlimited resources if the requests come in too fast. We can address that deficiency by changing Serve to gate the creation of the goroutines:

```go
func Serve(queue chan *Request) {
    for req := range queue {
        sem <- 1
        go func() {
            process(req)
            <-sem
        }()
    }
}
```
- Another approach that manages resources well is to start a fixed number of handle goroutines all reading from the request channel. The number of goroutines limits the number of simultaneous calls to process. This Serve function also accepts a channel on which it will be told to exit; after launching the goroutines it blocks receiving from that channel.

```go
func handle(queue chan *Request) {
    for r := range queue {
        process(r)
    }
}

func Serve(clientRequests chan *Request, quit chan bool) {
    // Start handlers
    for i := 0; i < MaxOutstanding; i++ {
        go handle(clientRequests)
    }
    <-quit  // Wait to be told to exit.
}

```
## Channels of Channels

- Channel is a first class value that can be allocated and passed around like any other value. A common usecase of this to implement safe parallel demultiplexing.

```go
type Request struct {
    args        []int
    f           func([]int) int
    resultChan  chan int
}
```

The client provides a function and its arguments, as well as a channel inside the request object on which to receive the answer.

```go
func sum(a []int) (s int) {
    for _, v := range a {
        s += v
    }
    return
}

request := &Request{[]int{3, 4, 5}, sum, make(chan int)}
// Send request
clientRequests <- request
// Wait for response.
fmt.Printf("answer: %d\n", <-request.resultChan)
```

## Parallelization

- Let's say we have to perform expensive operation vector of items and that the value of the operation on each item is independent. We can parallelize those ops

```go
type Vector []float64

// Apply the operation to v[i], v[i+1] ... up to v[n-1].
func (v Vector) DoSome(i, n int, u Vector, c chan int) {
    for ; i < n; i++ {
        v[i] += u.Op(v[i])
    }
    c <- 1    // signal that this piece is done
}
```

- We launch the pieces independently in a loop, one per CPU. They can complete in any order but it doesn't matter; we just count the completion signals by draining the channel after launching all the goroutines.

```go
const numCPU = 4 // number of CPU cores

func (v Vector) DoAll(u Vector) {
    c := make(chan int, numCPU)  // Buffering optional but sensible.
    for i := 0; i < numCPU; i++ {
        go v.DoSome(i*len(v)/numCPU, (i+1)*len(v)/numCPU, u, c)
    }
    // Drain the channel.
    for i := 0; i < numCPU; i++ {
        <-c    // wait for one task to complete
    }
    // All done.
}
```

- Rather than creating a constant number of CPU cores , We can use the ``runtime.NumCPU()`` to return the number of hardware cores in the machine and parallelize accordingly.

## Leaky buffer

- The client goroutine loops receiving data from some source, perhaps a network. 

- To avoid allocating and freeing buffers, it keeps a free list, and uses a buffered channel to represent it. If the channel is empty, a new buffer gets allocated. Once the message buffer is ready, it's sent to the server on ``serverChan``.

```go
var freeList = make(chan *Buffer, 100)
var serverChan = make(chan *Buffer)

func client() {
    for {
        var b *Buffer
        // Grab a buffer if available; allocate if not.
        select {
        case b = <-freeList:
            // Got one; nothing more to do.
        default:
            // None free, so allocate a new one.
            b = new(Buffer)
        }
        load(b)              // Read next message from the net.
        serverChan <- b      // Send to server.
    }
}
```

- The server loop receives each message from the client, processes it, and returns the buffer to the free list.

- The client attempts to retrieve a buffer from freeList; if none is available, it allocates a fresh one. The server's send to freeList puts b back on the free list unless the list is full, in which case the buffer is dropped on the floor to be reclaimed by the garbage collector. (The default clauses in the select statements execute when no other case is ready, meaning that the selects never block.) This implementation builds a leaky bucket free list in just a few lines, relying on the buffered channel and the garbage collector for bookkeeping.