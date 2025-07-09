# Slices - Deep Dive

- Slices in golang help us to work with sequence of data in a easy manner.

## Arrays

- The slice type is an abstraction that is built on top of the arrays.

- An array type element defines the length and element type. Eg ```[4]int``` represents array of type 4 integers.

- An array size is fixed and length of the array is part of its type. ```[4]int``` and ```[5]int``` are distinct incompatible types.

- Arrays can be indexed in the usual way, so the expression ```s[n]``` accesses the nth element, starting from zero.

```go
var a [4]int
a[0] = 1
i := a[0]
// i == 1
```

- Arrays do not need to be initialized explicitly; the zero value of an array is a ready-to-use array whose elements are themselves zeroed.


```go
// a[2] == 0, the zero value of the int type
```


- The in-memory represntation of ```[4]int``` is just four integer values laid sequentially

![Slice-array](https://go.dev/blog/slices-intro/slice-array.png)

- Go arrays are values, An array variable denotes an entire array, It's not pointer to the first element. This means when we assign or around a value we make a copy . (To avoid the copy you could pass a pointer to the array, but then that’s a pointer to an array, not an array.)

- Array is sort of struct but indexed rather than named fields a fixed size, One way to think about arrays is as a sort of struct but with indexed rather than named fields: a fixed-size composite value.

- We can either specify a size or make the compiler infer the size when initialising an array

```go
b := [...]string{"Penn", "Teller"} // Evlauated to 2 during runtime
b := [2]int{1,2}

```
---

## Slices

- Arrays are bit inflexible , So we don't see them often in the Gocode , Slices are built on top of array to provide great power and consistence.

- The type specification for a slice is []T, where T is the type of the elements of the slice. Unlike an array type, a slice type has no specified length.

- A slice literal is declared just like an array literal, except you leave out the element count:

```go
a := []string{"v","i","s","h"}
```

- A slice can be created with the built-in function called make, which has the signature,

```go
func make([]T, len, cap) []T
```

- Where ``T`` is type of element that to be created , The ``make`` function takes type , length and capacity (option), When called, make allocates an array and returns a slice that refers to that array.

```go
var s []byte
s = make([]byte, 5, 5)
// s == []byte{0, 0, 0, 0, 0}
```
- When the capacity is omitted , the capacity defaults to length

```go
s := make([]byte, 5)

len(s) == 5 // true
cap(s) == 5 // true
```

- The zero value of a slice is ``nil``. The len and cap functions will both return 0 for a nil slice.

- A slice can also be formed by ``slicing`` an existing slice or array . Slicing is done by ```b[1:4]```, creates elements from 1 through 3 index. 


```go
// b[:2] == []byte{'g', 'o'} // starting to 1
// b[2:] == []byte{'l', 'a', 'n', 'g'}  // starting from 2 to rest
// b[:] == b // true point to same array
```

In Go, **slice comparison** using `==` **is not allowed**, except when comparing to `nil`.

### Example:

```go
a := []int{1, 2, 3}
b := []int{1, 2, 3}

// fmt.Println(a == b) // ❌ compile error: invalid operation: a == b
```

---

### Why?

Because slices are **not directly comparable** — they are **structs internally** containing:

* a pointer to the underlying array,
* length,
* capacity.

So Go prevents direct comparison using `==` to avoid ambiguity (pointer vs value).

---

### What can you compare?

#### 1. Compare against `nil`:

```go
var x []int
fmt.Println(x == nil) // ✅ true
```

#### 2. Use `reflect.DeepEqual` (slower, deep comparison):

```go
import "reflect"

a := []int{1, 2, 3}
b := []int{1, 2, 3}
fmt.Println(reflect.DeepEqual(a, b)) // ✅ true
```

#### 3. Write manual comparison (faster than reflect):

```go
func equal(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
```

---

### Summary:

| Comparison          | Allowed? | Compares                      |
| ------------------- | -------- | ----------------------------- |
| `slice1 == slice2`  | ❌        | Compile-time error            |
| `slice1 == nil`     | ✅        | Pointer comparison            |
| `reflect.DeepEqual` | ✅        | Element-wise deep comparison  |
| Manual loop         | ✅        | Element-wise value comparison |

So slices are not directly comparable by value or memory address. You need to explicitly choose how you want to compare them.


## Internals

- Slice is a descriptor of the array segment , It consists of pointer to an array , length of the segment and its capacity(the max length of the segment)

![Slice Struct](https://go.dev/blog/slices-intro/slice-struct.png)

-  Our variable ``s``, created earlier by ``make([]byte, 5)``, is structured like this


![Slice-1](https://go.dev/blog/slices-intro/slice-1.png)

- The length is the number of elements , The capacity is the number of elements that the array can hold

- As we slice s, observe the changes in the slice data structure and their relation to the underlying array:

```go
s = s[2:4]
```

![Slicing](https://go.dev/blog/slices-intro/slice-2.png)

- Slicing doesn't copy the slices data , It creates a new slice value that points to the original array.

- This makes slice operations as efficient as manipulating array indices. 

-  Therefore, modifying the elements (not the slice itself) of a re-slice modifies the elements of the original slice:

```go
d := []byte{'r', 'o', 'a', 'd'}
e := d[2:]
// e == []byte{'a', 'd'}
e[1] = 'm'
// e == []byte{'a', 'm'}
// d == []byte{'r', 'o', 'a', 'm'}
```

- Earlier we sliced s to a length shorter than its capacity. We can grow s to its capacity by slicing it again:

![alt text](https://go.dev/blog/slices-intro/slice-3.png)

- A slice cannot be grown beyond its capacity. Attempting to do so will cause a runtime panic, just as when indexing outside the bounds of a slice or array. Similarly, slices cannot be re-sliced below zero to access earlier elements in the array.


## Growing slices (the copy and append functions)

- To increase the capacity of a slice one must create a new, larger slice and copy the contents of the original slice into it. 

```go
t := make([]byte, len(s), (cap(s)+1)*2) // +1 in case cap(s) == 0
for i := range s {
        t[i] = s[i]
}
s = t
```

- The looping piece of this common operation is made easier by the built-in copy function. As the name suggests, copy copies data from a source slice to a destination slice. It returns the number of elements copied.

```go
func copy(dst, src []T) int
```

- The copy fn supports copying between slices of different length, it will copy only up to the smaller number of element. In addition coopy can handle src and dest slices that share the same underlying array.

```go
t := make([]byte, len(s), (cap(s)+1)*2)
copy(t, s)
s = t
```
- To append an element to the end of an array can use the append fn

```go
func append(s []T, x ...T) []T
```

- The append function appends the elements x to the end of the slice s, and grows the slice if a greater capacity is needed.

```go
a := make([]int, 1)
// a == []int{0}
a = append(a, 1, 2, 3)
// a == []int{0, 1, 2, 3}
```

- Appending one slice to another

```go
a := []string{"John", "Paul"}
b := []string{"George", "Ringo", "Pete"}
a = append(a, b...) // equivalent to "append(a, b[0], b[1], b[2])"
// a == []string{"John", "Paul", "George", "Ringo", "Pete"}
```

- 