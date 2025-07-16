# Scheduler
---
## Blocking

- In Go all I/O is blocking, The Go ecosystem is written in such a way that write against a blocking interface and handle concurrency through goroutines and channels rathen than callbacks and futures.

- An example for this is net/http package which will create a goroutine to handle the request whenever HTTP request comes.

- We can write the handler in a straight forward manner , if a request comes do this and do that without worrying about the inner workings.

## What does Goruntime need with a scheduler?

- The POSIX thread API is an logical extension to the existing UNIX process model and as such threads get lot of control as same as the processes.

- Threads have their own signal mask , Can be put into cgroups , assigned CPU affinity and queried for resources that they can use.

- All these features may not be techincally needed in userspace , if we follow this model it will cause a lot of overhead when we want to spawn 100k gooroutines.

- OS can't make informed scheduling decisions based on the Go model, Eg: Go garbage collectors require that all threads are stopped when running a collection for the memory to be in a consistent state. This involves waiting for running thread to reach a point to know that they are in consistent state.

- When we have lots of threads scheduled in the wild chances is that we will be waiting for a very long time to resolve to an stable state.

## M:N scheduler

- There are 3 usual models for threading
    - ``N:1``: Where several userspace threads are run on one OS thread. This has advantage of quick context switch but it cannot make full use of the multicore processsors.
    - Another is ``1:1 ``where one thread of execution matches one OS thread. It takes advantage of all of the cores on the machine, but context switching is slow because it has to trap through the OS.

- Go tries best of both the worlds using the ``M:N`` scheduler. It schedules arbitary number of goroutines into arbitary number of OS threads. We get quick context switches and can make use of multiple cores in the system. The main advantage is the complexity that it adds to scheduler.

![alt text](https://morsmachine.dk/our-cast.jpg)

- The triangle represents an OS thread. It's the thread of execution managed by the OS and works pretty much like your standard POSIX thread. In the runtime code, it's called M for machine.

- The circle represents a goroutine. It includes the stack, the instruction pointer and other information important for scheduling goroutines, like any channel it might be blocked on. In the runtime code, it's called a G.

The rectangle represents a context for scheduling. You can look at it as a localized version of the scheduler which runs Go code on a single thread. It's the important part that lets us go from a N:1 scheduler to a M:N scheduler. In the runtime code, it's called P for processor.

![alt text](https://morsmachine.dk/in-motion.jpg)

