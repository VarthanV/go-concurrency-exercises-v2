# Blank identifier

- The blank identifier can be assigned or declared with any value of any type, with the value discarded harmlessly

-  It's a bit like writing to the Unix /dev/null file: it represents a write-only value to be used as a place-holder where a variable is needed but the actual value is irrelevant.

## Blank identifier in multiple assignment

- The use of a blank identifier in a for range loop is a special case of a general situation: multiple assignment.

- If an assignment requires multiple values on the left side, but one of the values will not be used by the program, a blank identifier on the left-hand-side of the assignment avoids the need to create a dummy variable and makes it clear that the value is to be discarded

```go
if _, err := os.Stat(path); os.IsNotExist(err) {
    fmt.Printf("%s does not exist\n", path)
}
```
- It is not recommended to discard error value inorder to ignore the error

```go
// Bad! This code will crash if path does not exist.
fi, _ := os.Stat(path)
if fi.IsDir() {
    fmt.Printf("%s is a directory\n", path)
}
```

# Unused imports and variables

- It is an error to import a package or declare a variable without using it.

- Unused imports bloat the program and slow compilation, where a variable that is initialized but it is not used, results in a wasted computation and indication of a larger bug.

- . When a program is under active development, however, unused imports and variables often arise and it can be annoying to delete them just to have the compilation proceed, only to have them be needed again later. The blank identifier provides a workaround.

- To silence complaints about the unused imports, use a blank identifier to refer to a symbol from the imported package. Similarly, assigning the unused variable fd to the blank identifier will silence the unused variable error. This version of the program does compile.


```go
package main

import (
    "fmt"
    "io"
    "log"
    "os"
)

var _ = fmt.Printf // For debugging; delete when done.
var _ io.Reader    // For debugging; delete when done.

func main() {
    fd, err := os.Open("test.go")
    if err != nil {
        log.Fatal(err)
    }
    // TODO: use fd.
    _ = fd
}

```

- By convention, the global declarations to silence import errors should come right after the imports and be commented, both to make them easy to find and as a reminder to clean things up later.

# Import for sideeffect

- It is useful to import a package only for its side effects.

- For example during its init function ``net/http/pprof`` package registers HTTP handlers to provide debugging information.

-  It has an exported API, but most clients need only the handler registration and access the data through a web page.

```go
import _ "net/http/pprof"
```

- This form of import makes clear that the package is being imported for its side effects, because there is no other possible use of the package: in this file, it doesn't have a name.


# Interface checks

- A type need not declare it implements an interface , Instead a type implements an interface just implementing the methods of the interface.

- In practice most of the interface conversions are static so they are checked during runtime.

- For example, passing an ``*os.File`` to a function expecting an ``io.Reader ``will not compile unless ``*os.File`` implements the io.Reader interface.

- One instance is in the encoding/json package, which defines a ``Marshaler`` interface. When the JSON encoder receives a value that implements that interface

```go
m, ok := val.(json.Marshaler)
```

- If it is necessary to check whether a type implements an interface and without actually using the interface itself,perhaps not part of an error check, use the blank identifier to ignore the same.

```go
if _, ok := val.(json.Marshaler); ok {
    fmt.Printf("value %v of type %T implements json.Marshaler\n", val, val)
}
```

- Don't do this for every type that satisfies an interface, though. By convention, such declarations are only used when there are no static conversions already present in the code, which is a rare event.
