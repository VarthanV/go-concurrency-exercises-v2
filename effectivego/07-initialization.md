# Initialization

- Initialization in Go is more powerful. Complex structures can be built during initialization and the ordering issues among initialized objects, even among different packages, are handled correctly.

## Constants

- Constants in Go are just that- constant. 

-  They are created a compile time even when defined as locals in functions and can only be **numbers,characters(runes),strings, or boolean**.

- Because of the compile time restrictions the expressions that they define must be constant expressions. 

- In Go, enumerated constants are created using the iota enumerator. Since iota can be part of an expression and expressions can be implicitly repeated, it is easy to build intricate sets of values.


```go
type ByteSize float64

const (
    _           = iota // ignore first value by assigning to blank identifier
    KB ByteSize = 1 << (10 * iota)
    MB
    GB
    TB
    PB
    EB
    ZB
    YB
)
```

### ✅ What is an `enum`?

An **`enum`** (short for **enumeration**) is a **user-defined data type** that consists of **named constant values**. It is used to represent a **fixed set of related values**, usually in a **readable** and **type-safe** way.

---

## 🔹 1. General Concept

Instead of using magic numbers (like `0`, `1`, `2`), `enum` gives meaningful **names** to constants.

---

### ✅ Example in **C / C++**:

```c
enum Color {
    RED,    // 0
    GREEN,  // 1
    BLUE    // 2
};

enum Color c = GREEN;
```

> Internally, `RED = 0`, `GREEN = 1`, `BLUE = 2`, unless you override them.

You can also assign custom values:

```c
enum Status {
    SUCCESS = 200,
    ERROR = 500
};
```

---

### ✅ Example in **Go** (using `iota`):

```go
type Status int

const (
    Pending Status = iota  // 0
    Approved               // 1
    Rejected               // 2
)
```

---

### ✅ Example in **Java**:

Java `enum` is more powerful (it creates a class-like type):

```java
enum Day {
    MONDAY, TUESDAY, WEDNESDAY
}

Day today = Day.MONDAY;
```

---

## 🔹 Why Use `enum`?

* Improves **readability**
* Enforces **valid values only**
* Reduces bugs from using arbitrary numbers
* Enables **switch/case** logic more cleanly

---

## 🔹 Enum vs Enumerator

| Term         | Meaning                                                                                  |
| ------------ | ---------------------------------------------------------------------------------------- |
| `enum`       | A **type** that holds a fixed set of named constants                                     |
| `enumerator` | A **single constant** inside the enum OR an **iterator** (depending on language context) |

---

Would you like a visual explanation or usage in your favorite language like Go or Python?

## Init function

- Each file can define its own nildaic ``init`` function to setup whatever state is required.

- Each file can have multiple init functions.

- init is called after all the variable declarations in the package have evaluated their initializers, and those are evaluated only after all the imported packages have been initialized.

- Besides initializations that cannot be expressed as declarations, a common use of init function is to verify the correctness of the program state before real execution begins

```go
func init() {
    if user == "" {
        log.Fatal("$USER not set")
    }
    if home == "" {
        home = "/home/" + user
    }
    if gopath == "" {
        gopath = home + "/go"
    }
    // gopath may be overridden by --gopath flag on command line.
    flag.StringVar(&gopath, "gopath", gopath, "override default GOPATH")
}
```