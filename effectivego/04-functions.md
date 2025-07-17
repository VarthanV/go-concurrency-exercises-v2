# Functions

- Go's feature is functions and methods can return multiple values.

- In Write method we can return the bytes written and error , There is possiblity that we have filled only partial bytes

```go
func (file *File) Write(b []byte) (n int, err error)
```

- It returns the number of bytes written and a non-nil error when ``n != len(b)``.

- A similar approach obviates the need to pass a pointer to a return value to simulate a reference parameter. Here's a simple-minded function to grab a number from a position in a byte slice, returning the number and the next position.

```go
func nextInt(b []byte, i int) (int, int) {
    for ; i < len(b) && !isDigit(b[i]); i++ {
    }
    x := 0
    for ; i < len(b) && isDigit(b[i]); i++ {
        x = x*10 + int(b[i]) - '0'
    }
    return x, i
}
```

- Can use to scan the numbers in input slice b like this

```go
  for i := 0; i < len(b); {
        x, i = nextInt(b, i)
        fmt.Println(x)
    }
```

# Named result parameters

- The return or result parameters of a Go function can be given names and used as regular variables, just like the incoming parameters.

- When named they are initialized to zero value of their types when the fn begins.

-  If the function executes a return statement with no arguments, the current values of the result parameters are used as the returned values.

- The names are not mandatory but they can make code shorter and clearer: they're documentation. If we name the results of nextInt it becomes obvious which returned int is which.

```go
func nextInt(b []byte, pos int) (value, nextPos int) 
```

- Because named results are initialized and tied to an unadorned return, they can simplify as well as clarify. Here's a version of io.ReadFull that uses them well:

```go
func ReadFull(r Reader, buf []byte) (n int, err error) {
    for len(buf) > 0 && err == nil {
        var nr int
        nr, err = r.Read(buf)
        n += nr
        buf = buf[nr:]
    }
    return
}
```

# Defer

- Go's defer statement schedules a fn call to be run immediately before the fn executing the defer returns.

- It's unusual but effective way to deal with releasing of resources regardless of fn succeeds or fails , like mutex lock unreleasing or closing a file

```go
// Contents returns the file's contents as a string.
func Contents(filename string) (string, error) {
    f, err := os.Open(filename)
    if err != nil {
        return "", err
    }
    defer f.Close()  // f.Close will run when we're finished.

    var result []byte
    buf := make([]byte, 100)
    for {
        n, err := f.Read(buf[0:])
        result = append(result, buf[0:n]...) // append is discussed later.
        if err != nil {
            if err == io.EOF {
                break
            }
            return "", err  // f will be closed if we return here.
        }
    }
    return string(result), nil // f will be closed if we return here.
}
```

- Deferring a fn call to ``Close`` has two advantages , it guarantees that we will never miss to close a file , closing of a file sits near open which is clever than placing it in the end of the file.

- The arguments to the deferred function (which include the receiver if the function is a method) are evaluated when the defer executes, not when the call executes.

```go
for i := 0; i < 5; i++ {
    defer fmt.Printf("%d ", i)
}
```
- Deferred functions are executed in ``LIFO`` order , so the above code will print ``4 3 2 0``.

- We can also write some utils to trace the function execution through out the program

```go
func trace(s string)   { fmt.Println("entering:", s) }
func untrace(s string) { fmt.Println("leaving:", s) }

// Use them like this:
func a() {
    trace("a")
    defer untrace("a")
    // do something....
}
```

```go
defer fmt.Println(x)
```

- The value of x is evaluated immediately at the point when the defer statement is encountered — not when the deferred function is actually executed.