# Synchornisation Effect

- A synchornisation effect mean an action like unlocking a mutex or closing a channel  creates a guaranteed ordering of memory access between two goroutines.

- One goroutine sees the up-to-date changes made by another goroutine.

- Without synchronization, goroutines can see stale or inconsistent values due to compiler optimizations, CPU caching, or out-of-order execution.

## Why it is important?

- Because goroutines run concurrently, one might modify memory while another is reading it. Without synchronization:
    - There’s no guarantee the reader sees the latest value.
    - You may get race conditions (nondeterministic or buggy behavior).


##  Examples of things with synchronization effect:

- Unlock → Lock (on a sync.Mutex)

- Send → Receive (on a channel)

- Close → Receive (zero value) (on a channel)

- RUnlock → Lock (on a sync.RWMutex)

- Successful TryLock()

All these ensure that  
    Goroutine B will see everything Goroutine A did before this synchronization point."


## What does not have a synchronization effect?

- A failed TryLock

- Reading a shared variable without any locks or channels

- Spawning a goroutine without any sync (like a waitgroup, channel, or mutex)


```go
var mu sync.Mutex
var x int

func writer() {
    x = 42          // write
    mu.Unlock()     // sync point
}

func reader() {
    mu.Lock()       // sync point
    fmt.Println(x)  // guaranteed to print 42
}
```

Synchronization effect = visibility guarantee between goroutines.

