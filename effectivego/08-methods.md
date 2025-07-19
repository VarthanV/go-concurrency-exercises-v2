# Methods

## Pointers vs Values

```go
type ByteSlice []byte

func (slice ByteSlice) Append(data []byte) []byte {
    // Body exactly the same as the Append function defined above.
}
```

- In this method we need to return a new byte to the function after appending and then need to modify the original byteslice in the caller.

- Instead we redefine the method to take a pointer to the ``ByteSlice`` as its receiver so it can overwrite the callers slice.

```go
func (p *ByteSlice) Append(data []byte) {
    slice := *p
    // Body as above, without the return.
    *p = slice
}
```

- In fact, we can do even better. If we modify our function so it looks like a standard Write method, like this,

```go
func (p *ByteSlice) Write(data []byte) (n int, err error) {
    slice := *p
    // Again as above.
    *p = slice
    return len(data), nil
}
```

- 