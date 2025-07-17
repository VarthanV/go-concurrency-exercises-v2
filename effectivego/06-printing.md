# Printing

- Formatted printing in Go uses a style similar to C's printf family but is richer and more general. T

- The functions live in the fmt package and have capitalized names: f``mt.Printf, fmt.Fprintf, fmt.Sprintf`` and so on

- The string functions (``Sprintf`` etc.) return a string rather than filling in a provided buffer.

```go
fmt.Printf("Hello %d\n", 23)
fmt.Fprint(os.Stdout, "Hello ", 23, "\n")
fmt.Println("Hello", 23)
fmt.Println(fmt.Sprint("Hello ", 23))
```

- The formatted print functions fmt.Fprint and friends take as a first argument any object that implements the io.Writer interface; the variables os.Stdout and os.Stderr are familiar instances.

- If we want the exact value that is represented to be presented can used the ``%v`` format specifier.

```go
fmt.Printf("%v\n", timeZone)  // or just fmt.Println(timeZone)
```

- For maps, Printf and friends ``sort the output lexicographically by key``.

- When printing a struct, the modified format ``%+v ``annotates the fields of the structure ``with their names,`` and for any value the alternate format ``%#v`` prints the ``value in full Go syntax``.

```go
type T struct {
    a int
    b float64
    c string
}
t := &T{ 7, -2.35, "abc\tdef" }
fmt.Printf("%v\n", t)
fmt.Printf("%+v\n", t)
fmt.Printf("%#v\n", t)
fmt.Printf("%#v\n", timeZone)
```
prints

```go
&{7 -2.35 abc   def}
&{a:7 b:-2.35 c:abc     def}
&main.T{a:7, b:-2.35, c:"abc\tdef"}
map[string]int{"CST":-21600, "EST":-18000, "MST":-25200, "PST":-28800, "UTC":0}
```

- The format ``%T`` prints type of a value

```go
fmt.Printf("%T\n", timeZone)
```

- If we want to control the default format for custom type, All needed is to define a method with a signature ``String() string`` 

```go
func (t *T) String() string {
    return fmt.Sprintf("%d/%g/%q", t.a, t.b, t.c)
}
```

Will print in the format

```go
7/-2.35/"abc\tdef"
```

- If need to print values of type T as well as pointer to T, the receiver of ``String`` must be of value type.

# Append 

- The append function returns a new slice.

```go
func append(slice []T, elements ...T) []T
```

- T is placeholder for an given type.

- New slice is returned when we append values to slice is because , the size of the underlying array may be declared differently if the array is resized this slice may be pointing to the new array.

