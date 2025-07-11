## Mutex

## sync.Mutex

- A basic mutual exclusion lock only one goroutine can hold it at a time, others block until it is locked.

- **Common usecase**: Protect read/write access to a shared variable or critical section.

```go
var mu sync.Mutex
var count int

func increment() {
	mu.Lock()
	count++
	mu.Unlock()
}

```

- **When to use**: Simple critical sections (read+write), No need for concurrent reads, Multiple writers or mixed access


## sync.RWMutex (Read-Write Mutex)

- A lock that distinguishes between readers and writers.

- Many goroutines can hold the lock for reading (``RLock``) at the same time.

- But only one can hold it for writing (``Lock``), and no reads are allowed then.

- **Common usecase**: Optimize performance for read-heavy workloads.

```go
var rw sync.RWMutex
var config map[string]string

func readConfig(key string) string {
	rw.RLock()
	defer rw.RUnlock()
	return config[key]
}

func writeConfig(key, value string) {
	rw.Lock()
	defer rw.Unlock()
	config[key] = value
}

```
- **When to use**

- Many reads, few writes

- Caching, configuration maps, lookup tables


## sync.Once

- Ensures that a function only runs once, even across multiple goroutines.

- **Common usecase**: Lazy initialization, Singleton setup

```go
var once sync.Once
func initConfig() {
	once.Do(func() {
		fmt.Println("Initializing config...")
	})
}
```

- **When to use**: Initialize shared resources only once, Safe singletons, Avoid race conditions in setup code.


## sync.TryLock()

- Attempts to lock a ``sync.Mutex ``without blocking. Returns ``true`` if successful, ``false`` otherwise.

- Try doing work only if the resource is not locked, else skip or defer.


```go
var mu sync.Mutex

if mu.TryLock() {
	defer mu.Unlock()
	fmt.Println("Acquired lock, doing work")
} else {
	fmt.Println("Could not get lock, skipping")
}

```

- **When to use**: Non-blocking access, Optional work, Deadlock avoidance in special cases


| Type        | Description               | Concurrency             | Use Case                               |
| ----------- | ------------------------- | ----------------------- | -------------------------------------- |
| `Mutex`     | Exclusive lock            | ❌ No reads while locked | Simple read/write protection           |
| `RWMutex`   | Read/Write lock           | ✅ Multiple readers      | Read-heavy shared data access          |
| `Once`      | Run a function only once  | ✅ Safe for all          | Lazy init, setup code                  |
| `TryLock()` | Non-blocking lock attempt | ⚠️ Careful use          | Skip if busy, optional fast-path logic |


## Once

- The sync package provides a safe mechanism for initialization in the presence of multiple goroutines through the use of the ``Once`` type.

- Multiple threads can execute ``once.Do(f) ``for a particular f, but only one will run ``f()``, and the other calls block until f() has returned.

- The completion of a single call of f() from once.Do(f) is synchronized before the return of any call of once.Do(f).


```go
var a string
var once sync.Once

func setup() {
	a = "hello, world"
}

func doprint() {
	once.Do(setup)
	print(a)
}

func twoprint() {
	go doprint()
	go doprint()
}
```

calling twoprint will call setup exactly once. The setup function will complete before either call of print. The result will be that "hello, world" will be printed twice.

## Atomic Values

- The APIs in the sync/atomic package are collectively “atomic operations” that can be used to synchronize the execution of different goroutines.

-  If the effect of an atomic operation A is observed by atomic operation B, then A is synchronized before B. All the atomic operations executed in a program behave as though executed in some sequentially consistent order.

## Finalizers

- The runtime package provides a SetFinalizer function that adds a finalizer to be called when a particular object is no longer reachable by the program. A call to SetFinalizer(x, f) is synchronized before the finalization call f(x).

