**Synchronization** in multithreading or goroutines is the coordination of threads/goroutines to ensure:

1. **Correctness** — Data is not corrupted or inconsistently read/written.
2. **Ordering** — Certain operations happen in the right sequence.
3. **Mutual exclusion** — Only one thread accesses critical sections at a time.

---

### 🔹 Why Synchronization is Needed

When multiple threads/goroutines access **shared resources (e.g., variables, files, memory)**, there's a risk of:

* **Race conditions** – Two threads/goroutines access a variable simultaneously.
* **Data inconsistency** – One writes while another reads.
* **Deadlocks/livelocks** – Threads block each other indefinitely.

---

### 🔹 Common Synchronization Methods

#### ✅ In Golang

1. **Mutex (`sync.Mutex`)**
   Ensures only one goroutine accesses the critical section.

   ```go
   var mu sync.Mutex
   mu.Lock()
   // critical section
   mu.Unlock()
   ```

2. **WaitGroup (`sync.WaitGroup`)**
   Waits for multiple goroutines to finish.

   ```go
   var wg sync.WaitGroup
   wg.Add(1)
   go func() {
       defer wg.Done()
       // do work
   }()
   wg.Wait()
   ```

3. **Channels**
   Built-in tool for goroutine synchronization via communication.

   ```go
   ch := make(chan int)
   go func() {
       ch <- 42  // send synchronizes
   }()
   val := <-ch  // receive synchronizes
   ```

4. **`sync.Once`, `sync.Cond`, `atomic`** — Other advanced tools for controlled synchronization.

---

### 🔹 Memory Synchronization

* Beyond just timing, **memory synchronization** ensures changes made by one goroutine are visible to others.
* Go provides a **happens-before** guarantee: certain actions (like `go` statement, `channel send/receive`, `Lock/Unlock`) form memory synchronization points.

> Example: A `go` statement synchronizes **before** the new goroutine starts.

---

### 🔹 In Simple Terms

**Synchronization** is like:

* A traffic signal ensuring cars (threads) don’t crash.
* A relay race baton — one can't run until the baton is passed.

---