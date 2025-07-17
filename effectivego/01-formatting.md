# Formatting

- Formatting issues are the most contentious but least consequential,It is best to leave the formatting decisions ``gofmt``.

- ``gofmt`` formats files at a package level.

- Example lets say this is the below unformatted code

```go
type T struct {
    name string // name of the object
    value int // its value
}
```

- gofmt automatically formats it to this

```go
type T struct {
    name    string // name of the object
    value   int    // its value
}
```

# Commentary

- Go provides C style ``/**/`` block comments and C++ style ``//`` line comments, Line comments are the norm, Block comments appear mostly as package level comments but useful within an expression or to disable large swath of code.

- Comments that appear before top-level declarations, with no intervening newlines, are considered to document the declaration itself. These “doc comments” are the primary documentation for a given Go package or command.

Eg:

```go
// Package path implements utility routines for manipulating slash-separated
// paths.
//
// The path package should only be used for paths separated by forward
// slashes, such as the paths in URLs. This package does not deal with
// Windows paths with drive letters or backslashes; to manipulate
// operating system paths, use the [path/filepath] package.
package path
```


