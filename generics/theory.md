# Generics

- Generics allows to write functions, types (structs, interfaces), and methods that operate on types that are specified later.

- Instead of writing separate functions for int, float64, and string for a common operation like "finding the minimum,"  can write one generic function that works for any type that satisfies certain requirements.

## Before Generics (The Problem Generics Solve):

- Imagine you want to write a Min function for different numeric types:

```go
func MinInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func MinFloat64(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
// ... and so on for int32, int64, float32, etc.
```

- This leads to repetitive code (boilerplate). Alternatively,  might use ``interface{}`` and type assertions, which sacrifices type safety and introduces runtime overhead.

## Key concepts of Go generics

## Type parameters
---
- Generics use type parameters to represent types that will be determined at the time of usage.

- Type parameters are enclosed in square brackets [] after the function name or type name.

```go
    // T is a type parameter
    func PrintValue[T any](value T) {
        fmt.Println(value)
    }
```

## Type constraints

- Type parameters need constraints to define what operations are allowed on the generic type. Constraints are simply interfaces.

- They specify the set of types that can be used as type arguments.

- ``any``: The most permissive constraint, ``equivalent to interface{}``, means any type can be used. It allows you to store and pass values of any type, but you can only perform operations valid for any (``e.g., assignment, comparison with == or != if the underlying type is comparable``).

- ``comparable``: A built-in constraint for types that can be compared using == and !=. This includes booleans, numbers, strings, pointers, channels, arrays of comparable types, and structs whose fields are all comparable types.


- ``Custom Interfaces as Constraints``: Can define your own interfaces to specify required methods. Any type that implements these methods will satisfy the constraint.


- ``Union of Types (|)``: Can specify a set of concrete types or underlying types using the | operator. 
    - ``int | float64``: Allows either int or float64.

    - ``~int | ~float64``: The ~ (tilde) symbol denotes the underlying type. This means it accepts int, float64, and any named type whose underlying type is int or float64 (e.g., type MyInt int).


## Generic functions

```go
func FunctionName[TypeParameters](parameters ParameterType) ReturnType {
    // ...
}
```

## Generic min function

```go
package main

import (
	"fmt"
	"golang.org/x/exp/constraints" // You might need to `go get golang.org/x/exp`
)

// Min returns the smaller of two values.
// T is constrained to be an Ordered type (e.g., int, float64, string).
func Min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func main() {
	fmt.Println("Min(5, 10):", Min(5, 10))         // int
	fmt.Println("Min(3.14, 2.71):", Min(3.14, 2.71)) // float64
	fmt.Println("Min(\"apple\", \"banana\"):", Min("apple", "banana")) // string

	type MyInt int
	var mi1 MyInt = 7
	var mi2 MyInt = 3
	fmt.Println("Min(MyInt(7), MyInt(3)):", Min(mi1, mi2)) // Custom type with underlying int
}
```
## Type Inference:

- Notice that in the main function, we don't explicitly specify ``Min[int](5, 10)``. The Go compiler can infer the type parameters from the arguments passed.

## Generic Types (Structs and interfaces)

## Generic Structs:
You can define structs that hold values of a generic type.

```go
package main

import "fmt"

// Stack is a generic stack data structure.
type Stack[T any] struct {
	elements []T
}

// Push adds an element to the stack.
func (s *Stack[T]) Push(item T) {
	s.elements = append(s.elements, item)
}

// Pop removes and returns the top element from the stack.
func (s *Stack[T]) Pop() (T, bool) {
	if len(s.elements) == 0 {
		var zero T // Return zero value for the type T
		return zero, false
	}
	lastIndex := len(s.elements) - 1
	item := s.elements[lastIndex]
	s.elements = s.elements[:lastIndex]
	return item, true
}

// IsEmpty checks if the stack is empty.
func (s *Stack[T]) IsEmpty() bool {
	return len(s.elements) == 0
}

func main() {
	// Create a stack of integers
	intStack := Stack[int]{}
	intStack.Push(10)
	intStack.Push(20)
	val, ok := intStack.Pop()
	fmt.Println("Popped from intStack:", val, ok) // Output: Popped from intStack: 20 true

	// Create a stack of strings
	stringStack := Stack[string]{}
	stringStack.Push("hello")
	stringStack.Push("world")
	sVal, sOk := stringStack.Pop()
	fmt.Println("Popped from stringStack:", sVal, sOk) // Output: Popped from stringStack: world true
}
```

## Generic Interfaces

- Generic interfaces define a set of methods that involve type parameters.

```go
package main

import "fmt"

// Processor is a generic interface for processing elements of type T.
type Processor[T any] interface {
	Process(item T) T
	Describe() string
}

// IntProcessor implements Processor for int.
type IntProcessor struct{}

func (ip IntProcessor) Process(item int) int {
	return item * 2
}

func (ip IntProcessor) Describe() string {
	return "Doubles an integer"
}

// StringProcessor implements Processor for string.
type StringProcessor struct{}

func (sp StringProcessor) Process(item string) string {
	return item + " (processed)"
}

func (sp StringProcessor) Describe() string {
	return "Appends '(processed)' to a string"
}

func PerformProcessing[T any](p Processor[T], data T) {
	fmt.Printf("Processor: %s, Input: %v, Output: %v\n", p.Describe(), data, p.Process(data))
}

func main() {
	intProc := IntProcessor{}
	PerformProcessing(intProc, 5) // Output: Processor: Doubles an integer, Input: 5, Output: 10

	stringProc := StringProcessor{}
	PerformProcessing(stringProc, "example") // Output: Processor: Appends '(processed)' to a string, Input: example, Output: example (processed)
}
```
## When to Use Generics (Best Practices)


- **Code Reusability:** When you find yourself writing almost identical code for different types (e.g., MinInt, MinFloat), generics are a strong candidate.


- **Type Safety**: Generics maintain Go's strong static typing at compile-time, preventing runtime errors that could occur with ``interface{}`` and type assertions.

- **Data Structures**: Implementing common data structures like stacks, queues, linked lists, trees, or generic collections.

- **Algorithms**: Writing algorithms that operate independently of the element type (e.g., sorting, filtering, mapping over slices).

- **Avoiding Reflection/Type Assertions**: If  goal is to handle different types dynamically,and are currently using reflection or extensive type assertions, generics might offer a more performant and type-safe solution.

## When not to use Generics

- **Simple Type-Specific Functions**: If a function inherently only makes sense for one or a very few specific types (e.g., a function dealing with HTTP requests), don't force generics.

- **Method Implementations Differ**: If the core logic of  operations would be significantly different for each type, interfaces with separate implementations are still the way to go.  Don't use generics if you'd end up with a huge switch statement based on type inside your generic function.

- **Over-Generalization**: Don't use generics just because you can. Sometimes, a well-defined interface or even a bit of duplication is clearer and more maintainable than an overly complex generic solution.

- **Performance-Critical Code (initial thought)**: While generics can improve performance by avoiding reflection, in some extremely performance-critical scenarios, the compiler might not optimize generic code as aggressively as concrete code. Always benchmark if performance is a critical concern. (Note: Go's compiler is constantly improving here.)

In Go generics, `constraints.Ordered` and `constraints.Comparable` are two useful **type sets** defined in the `golang.org/x/exp/constraints` package (or conceptually in the standard library as of Go 1.18+ for built-in types).

Here’s the difference between them:

---

### 🟦 `constraints.Comparable`

* This constraint allows any type that can be used with the **`==` and `!=`** operators.
* Basically, this is any type that is **comparable** in Go (can be keys in a map).
* **Includes**:

  * All basic types: `int`, `float64`, `string`, `bool`, `uintptr`
  * Structs (if all fields are comparable)
  * Arrays (if element type is comparable)
* **Does NOT include**:

  * Slices
  * Maps
  * Functions

#### Use case:

```go
func IndexOf[T comparable](slice []T, target T) int {
	for i, v := range slice {
		if v == target {
			return i
		}
	}
	return -1
}
```

---

### 🟨 `constraints.Ordered`

* This constraint allows types that support **ordering operators**: `<`, `<=`, `>`, `>=`
* Only includes basic **ordered types**:

  * `integers`: `int`, `int32`, `uint64`, etc.
  * `floats`: `float32`, `float64`
  * `string`
* **Does NOT include**:

  * `bool`
  * Structs, arrays, slices, maps

#### Use case:

```go
func Max[T constraints.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}
```

---

### Summary Table

| Constraint               | Supports `==`, `!=` | Supports `<`, `>` | Includes Structs/Arrays | Includes Strings | Includes Bools |
| ------------------------ | ------------------- | ----------------- | ----------------------- | ---------------- | -------------- |
| `constraints.Comparable` | ✅                   | ❌                 | ✅ (if fields are)       | ✅                | ✅              |
| `constraints.Ordered`    | ✅                   | ✅                 | ❌                       | ✅                | ❌              |

---

### TL;DR:

* Use `**Comparable**` when you need equality checks (e.g., searching).
* Use `**Ordered**` when you need to sort or compare greater/less than.

Let me know if you want a quick demo or real-world example with generics.

## Custom type constraints

- Can create custom type constraints by defining an interface that includes the methods that the type parameter must implement, or by specifying a union of types.

```go
type Stringer interface {
    String() string
}

func PrintAndDescribe[T Stringer](item T) {
    fmt.Println("Value:", item.String())
}
```

```go
type Numeric interface {
    int | int8 | int16 | int32 | int64 |
    uint | uint8 | uint16 | uint32 | uint64 | uintptr |
    float32 | float64 |
    ~int | ~float64 // Using ~ for underlying types
}

func Sum[T Numeric](slice []T) T {
    var total T
    for _, v := range slice {
        total += v
    }
    return total
}
```

## Interfaces vs Generics

- **Prefer Interfaces**: When  need to define a behavior (a set of methods) that different types can implement. If the operations on your types are polymorphic (i.e., different types perform the operation differently), interfaces are the correct choice. Interfaces focus on "what a type can do."

- **Prefer Generics**:  When the logic of your code is largely the same, but it needs to operate on different types that share a common structure or set of allowed operations (like arithmetic for numbers, or comparisons). Generics focus on "what types can be used."

 - If you're building a ``Shape`` interface with ``Area() ``and ``Perimeter()`` methods (because ``Circle`` and ``Square`` calculate them differently), use interfaces. If you're building a ``Filter`` function that works the same way regardless of whether it's filtering a slice of ints or strings, use generics.

