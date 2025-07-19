# Interfaces and other types

## Interfaces

- Interface in Go provide us a way to define the behaviour of the object, If something can do this then it can be used here.


## Interface Conversions and type assertions

- Type switches are form of conversion; they take an interface and for each case in switch ,convert it to the type of that case.

```go
type Stringer interface {
    String() string
}

var value interface{} // Value provided by caller.
switch str := value.(type) {
case string:
    return str
case Stringer:
    return str.String()
}
```

- If we care about only one type , onecase type switch would do , but so would a ``type assertion ``.

- A type assertion takes  an interface value and extracts value from it a value of specified explicit type. 

-  The syntax borrows from the clause opening a type switch, but with an explicit type rather than the type keyword.

```go
value.(typeName)
```
- The result is the new value with type statictype typeName. 

- That type must either be the concrete type held by the interface, or a second interface type that the value can be converted to.

```go
str := value.(string)
```

- But it turns out the value doesn't contain string , then the program will crash with a run-time error. 

- To guard against that, use the "comma, ok" idiom to test, safely, whether the value is a string

```go
str, ok := value.(string)
if ok {
    fmt.Printf("string value is: %q\n", str)
} else {
    fmt.Printf("value is not a string\n")
}
```

- If the type assertion fails, str will still exist and be of type string, but it will have the zero value, an empty string.

```go
if str, ok := value.(string); ok {
    return str
} else if str, ok := value.(Stringer); ok {
    return str.String()
}
```

## Interface and methods

- Since almost anything can have methods attached , almost anything can satisfy an interface.

- One example is the http package,which defines the ``Handler`` interface.

- Any object that implements Handler can serve HTTP requests.

```go
type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}
```
