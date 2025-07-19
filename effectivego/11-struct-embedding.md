# Embedding

- Go doesn't have the typical notion of subclassing , but it does have the ability to borrow pieces of implementation by embedding types within a struct or interface.

- Interface embedding is very simple

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}
```

- We can combine two interfaces by embedding two different interfaces into single interface

```go
// ReadWriter is the interface that combines the Reader and Writer interfaces.
type ReadWriter interface {
    Reader
    Writer
}
```

- Interfaces can be embedded only within interfaces

- The embedded elements are pointers to structs and of course must be initialized to point to valid structs before they can be used. The ReadWriter struct could be written as

```go
type ReadWriter struct {
    reader *Reader
    writer *Writer
}
```

- To promote the methods of the interface that reader implements need method forwarding like this

```go
func (rw *ReadWriter) Read(p []byte) (n int, err error) {
    return rw.reader.Read(p)
}
```

- By embedding the struct directly we avoid book keeping, The methods  of embedded types come along free which means the ``bufio.ReadWriter`` not only has the methods of bufio.Reader and bufio.Writer, it also satisfies all three interfaces: io.Reader, io.Writer, and io.ReadWriter.

- When we embed a type the methods of type become methods of the outer type,but when they are invoked the receiver of the method is inner type,not the outer one.

- When the Read method of a bufio.ReadWriter is invoked, it has exactly the same effect as the forwarding method written out above; the receiver is the reader field of the ReadWriter, not the ReadWriter itself.

- We can embed a field directly alongside a regular name field

```go
type Job struct {
    Command string
    *log.Logger
}
```

- The ``Job`` type now has  ``Println,Printf`` and other methods of ``*log.Logger``  We could have given the Logger a field name, of course, but it's not necessary to do so. And now, once initialized, we can log to the Job:

```go
job.Println("starting now...")
```

- The Logger is a regular field of the Job struct, so we can initialize it in the usual way inside the constructor for Job, like this,

```go
func NewJob(command string, logger *log.Logger) *Job {
    return &Job{command, logger}
}
```
or with a composite literal

```go
func (job *Job) Printf(format string, args ...interface{}) {
    job.Logger.Printf("%q: %s", job.Command, fmt.Sprintf(format, args...))
}
```

- 