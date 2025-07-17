# Control Structures

- There is ``no do or while`` loop , for switch is more flexible

- if and switch accept an ``optional initialization statement`` ,  like that for , break and continue statements take an optional label to identify what to break or continue.

- for; break and continue statements take an optional label to identify what to break or continue.


## If

- In Go a single if looks like this

```go
if x > y {
    return y
}
```

- Since if and switch accept an initialization statement, it's common to see one used to set up a local variable.

```go
if err := file.Chmod(0664); err != nil {
    log.Print(err)
    return err
}
```

- If an ``if`` statement ends with break, or return the else is omitted

```go
    f, err := os.Open(name)
    if err != nil {
        return err
    }
    codeUsing(f)
```

## Redeclaration and Assignment

```go
f, err := os.Open(name)
```

- The above statement declares two variables f and err, after few lines later ```f.Stat``` reads

```go
d, err := f.Stat()
```

- This seems like d and err are redeclared, The err appears in both statements, The duplication is legal , err is ``declared`` by the first statement , but ``reassigned`` in the second statement.

- Which means ``f.Stat`` uses the existing err variable declared and just gives a new value to the err declared.

- In a := declaration a variable v may appear even if it has already been declared, provided:
    - this declaration is in the same scope as the existing declaration of v (if v is already declared in an outer scope, the declaration will create a new variable §),
    - the corresponding value in the initialization is assignable to v, and
    - there is at least one other variable that is created by the declaration.

- This unusual property is pure pragmatism, making it easy to use a single err value.

## For

- Go's for loop is similar to but not same as the C's , It unifies for and while there is no do-while. There are three forms only one of which has semicolons

```go
// Like a C for
for init; condition; post { }

// Like a C while
for condition { }

// Like a C for(;;)
for { }
```

- Short declarations make it easy to declare the index variable in the right loop

```go
sum := 0
for i := 0; i < 10; i++ {
    sum += i
}
```

- If we are looping over an array, slice, string, or map, or reading from a channel, a range clause can manage the loop.

```go
for key, value := range oldMap {
    newMap[key] = value
}
```

- If we only need the first item in the range (key or index) can drop the second

```go
for key := range m {
    if key.expired() {
        delete(m, key)
    }
}
```
- If  only need the second item in the range (the value), use the blank identifier, an underscore, to discard the first:

```go
sum := 0
for _, value := range array {
    sum += value
}
```

- For strings , the range does more work by breaking out individual Unicode code point by parsing the UTF-8 , Erroneous encoding consume one byte and produce the replacement rune ``U+FFFD``.

```go
for pos, char := range "日本\x80語" { // \x80 is an illegal UTF-8 encoding
    fmt.Printf("character %#U starts at byte position %d\n", char, pos)
}

```

prints

```go
character U+65E5 '日' starts at byte position 0
character U+672C '本' starts at byte position 3
character U+FFFD '�' starts at byte position 6
character U+8A9E '語' starts at byte position 7
```

- Finally, Go has no comma operator and ++ and -- are statements not expressions. Thus if  want to run multiple variables in a for  should use parallel assignment (although that precludes ++ and --).

```go
// Reverse a
for i, j := 0, len(a)-1; i < j; i, j = i+1, j-1 {
    a[i], a[j] = a[j], a[i]
}

```

## Switch

- Switch is more general than C's.

- The expressions need not be constants or even integers, the cases are evaluated top to bottom until a match is found, and if the switch has no expression it switches on true.

- It is idiomatic to right a multiple if else chain into switch

```go
func unhex(c byte) byte {
    switch {
    case '0' <= c && c <= '9':
        return c - '0'
    case 'a' <= c && c <= 'f':
        return c - 'a' + 10
    case 'A' <= c && c <= 'F':
        return c - 'A' + 10
    }
    retu
```

- There is no automatic fallthrough but cases can be presented in comma seperated lists

```go
func shouldEscape(c byte) bool {
    switch c {
    case ' ', '?', '&', '=', '#', '+', '%':
        return true
    }
    return false
}
```

- Although they are not commonly used as in Go, break statements can be used to terminate a switch early. 

```go
Loop:
    for n := 0; n < len(src); n += size {
        switch {
        case src[n] < sizeOne:
            if validateOnly {
                break
            }
            size = 1
            update(src[n])

        case src[n] < sizeTwo:
            if n+1 >= len(src) {
                err = errShortInput
                break Loop
            }
            if validateOnly {
                break
            }
            size = 2
            update(src[n] + src[n+1]<<shift)
        }
    }
```

- The continue statement also accepts an optional label but it applies only to loops.

## Type Switch

- A switch can also be used to discover the dynamic type of an interface variable.

- Such a type switch uses the syntax of a type assertion with the keyword type inside the parentheses.

- If the switch declares a variable in the expression, the variable will have the corresponding type in each clause

- It's also idiomatic to reuse the name in such cases, in effect declaring a new variable with the same name but a different type in each case.

```go
var t interface{}
t = functionOfSomeType()
switch t := t.(type) {
default:
    fmt.Printf("unexpected type %T\n", t)     // %T prints whatever type t has
case bool:
    fmt.Printf("boolean %t\n", t)             // t has type bool
case int:
    fmt.Printf("integer %d\n", t)             // t has type int
case *bool:
    fmt.Printf("pointer to boolean %t\n", *t) // t has type *bool
case *int:
    fmt.Printf("pointer to integer %d\n", *t) // t has type *int
}
```