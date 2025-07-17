# Allocating with new

- Go has two allocation primitives the builtin function ``new`` and ``make``.

## new

- It is a builtin function that allocates , but unlike just initializing the value it also populates that memory with the zero values of the respective type`` T``.

- new(T) allocates zeroed storage for a new item of type T and returns its address, a value of type ``*T``.

- It returns a pointer to a newly allocated zero value of type T.

- Since the memory returned by new is zeroed, We can use it without further initialization.

- This means a user of the data structure can create one with new and get right to work.

- sync.Mutex does not have an explicit constructor or Init method. Instead, the zero value for a sync.Mutex is defined to be an unlocked mutex.

```go
type SyncedBuffer struct {
    lock    sync.Mutex
    buffer  bytes.Buffer
}
```

- Values of type ``SyncedBuffer`` are also ready to use immediately upon allocation or just declaration.

- Both v and p works immediately without further arrangement

```go
p := new(SyncedBuffer)  // type *SyncedBuffer
var v SyncedBuffer      // type  SyncedBuffer
```
# Constructors and composite literals

- Sometimes the zero value isn't good enough and an initializing constructor is necessary, as in this example derived from package os.

```go
func NewFile(fd int, name string) *File {
    if fd < 0 {
        return nil
    }
    f := new(File)
    f.fd = fd
    f.name = name
    f.dirinfo = nil
    f.nepipe = 0
    return f
}
```
- There's a lot of boilerplate in there. We can simplify it using a composite literal, which is an expression that creates a new instance each time it is evaluated.

```go
func NewFile(fd int, name string) *File {
    if fd < 0 {
        return nil
    }
    f := File{fd, name, nil, 0}
    return &f
}

```

- There's a lot of boilerplate in there. We can simplify it using a composite literal, which is an expression that creates a new instance each time it is evaluated.

```go
func NewFile(fd int, name string) *File {
    if fd < 0 {
        return nil
    }
    f := File{fd, name, nil, 0}
    return &f
}
```

## make

- The built-in function make(T, args) serves a purpose different from new(T). It creates **slices, maps, and channels only**, and it returns an initialized (not zeroed) value of type T (not *T).

- The reason for the distinction is that these three types represent, under the covers, references to data structures that must be initialized before use.

- A slice for example is a three item descriptor containing a pointer to the data (inside an array), and the capacity, and until those items are initialized, the slice is nil

- For slices, maps, and channels, make initializes the internal data structure and prepares the value for use.

```go
make([]int, 10, 100)
```

- Allocates an array of 100 ints and then creates a slice structure with length 10 and a capacity of 100 pointing at the first 10 elements of the array.

- When making a slice, the capacity can be omitted.

- Where as the new([]int) returns  a pointer to newly allocated zeroed slice structure that is a pointer to nil slice value.

```go
var p *[]int = new([]int)       // allocates slice structure; *p == nil; rarely useful
var v  []int = make([]int, 100) // the slice v now refers to a new array of 100 ints

// Unnecessarily complex:
var p *[]int = new([]int)
*p = make([]int, 100, 100)

// Idiomatic:
v := make([]int, 100)
```

- Remember that make applies only to maps, slices and channels and does not return a pointer. To obtain an explicit pointer allocate with new or take the address of a variable explicitly.

# Arrays

- Arrays are useful when planning detailed layout of memory and can sometimes help avoid allocation, they are primary building block of slices.

- Arrays are values , Assigning one array to another copies all elements.

- If we pass an array to a fn , The fn makes a copy of the array not a pointer to it.

- The size of an array is part of its type. The types ``[10]int ``and ``[20]int`` are distinct.

- Beter to use Go slices if we don't have very specific usecase to use arrays.

# Slices

- Slices wraps array to give more general powerful and convenient interface to sequence of data .

- Except for items with explicit dimension such as transformation matrices, most array programming in Go is done with slices rather than simple arrays.

- Slices holds refernece to an underlying array , if we assign one slice to another both refer to the same array.

- If a function takes a slice arguments , changes it the changes made inside the fn to the element will be visible to the caller, analogous of passing a pointer to an underlying array.

-  A Read function can therefore accept a slice argument rather than a pointer and a count. the length within the slice sets an upper limit of how much data to read.

```go
func (f *File) Read(buf []byte) (n int, err error)
```

- The method returns the number of bytes read and an error value, if any. To read into the first 32 bytes of a larger buffer buf, slice (here used as a verb) the buffer.

```go
    n, err := f.Read(buf[0:32])
```

- The capacity of a slice can be accessed by using builtin function ``cap``.

- The length of a slice may be changed as long as it still fits within the limits of the underlying array; just assign it to a slice of itself. 

- We must return the slice afterwards because, although ``Append`` can modify the elements of slice, the slice itself (the run-time data structure holding the pointer, length, and capacity) is passed by value.

# Two dimensional Slices

- Go's arrays and slices are one-dimensional. To create the equivalent of a 2D array or slice, it is necessary to define an array-of-arrays or slice-of-slices, like this:

```go
type Transform [3][3]float64  // A 3x3 array, really an array of arrays.
type LinesOfText [][]byte     // A slice of byte slices.
```

- Slices are variable-length, it is possible to have each inner slice be a different length. That can be a common situation, as in our LinesOfText example: each line has an independent length.

```go
text := LinesOfText{
    []byte("Now is the time"),
    []byte("for all good gophers"),
    []byte("to bring some fun to the party."),
}
```

## Maps

- Maps are convenient datastructure that associates value of one type (key) to value of other type(element or value).

- The key can be of any type for which the equality operator (==) is defined like integers,string,floating  points and complex numbers, strings, pointers, interfaces (as long as the dynamic type supports equality), structs and arrays

- Slices cannot be used as map keys since equality is not defined there.

- Like maps slices hold reference to an underlying data structure.

- If we pass a function to map and changes are made inside the function to map , the changes are visible in the caller.

- Maps can be constructed using the usual composite literal syntax with colon-separated key-value pairs, so it's easy to build them during initialization.

```go
var timeZone = map[string]int{
    "UTC":  0*60*60,
    "EST": -5*60*60,
    "CST": -6*60*60,
    "MST": -7*60*60,
    "PST": -8*60*60,
}
```

- Assigning and fetching map values looks syntactically just like doing the same for arrays and slices except that the index doesn't need to be an integer.

```go
offset := timeZone["EST"]
```

- An attempt to fetch a map value with type that doesn't exist will return the zero value of the map,If map contains an integer looking up an non-existential key will return 0. A set can be implemented as map with a value type bool.

-  Set the map entry to true to put the value in the set, and then test it by simple indexing.

```go
attended := map[string]bool{
    "Ann": true,
    "Joe": true,
    ...
}

if attended[person] { // will be false if person is not in the map
    fmt.Println(person, "was at the meeting")
}
```

- Sometimes we need to distinguish a missing entry from zero value,Is there an entry for "UTC" or is that 0 because it's not in the map at all? You can discriminate with a form of multiple assignment.

```go
func offset(tz string) int {
    if seconds, ok := timeZone[tz]; ok {
        return seconds
    }
    log.Println("unknown time zone:", tz)
    return 0
}
```

- To test the presence without worrying about the value can ignore the return value and just make decision based on the ok param

```go
_, present := timeZone[tz]
```

- To delete a map entry can use the built-in delete function

```go
delete(timeZone, "PDT")  // Now on Standard Time
```

