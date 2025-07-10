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

