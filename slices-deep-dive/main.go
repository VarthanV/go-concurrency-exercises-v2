package main

import "fmt"

func simpleArray() {
	var a [5]int
	a[0] = 1
	fmt.Println(a[0]) // Output: 1
	fmt.Println(a)    // Output: [1 0 0 0 0] , values are already initialized to 0 in
	// the capacity

}

func customSliceDriver() {
	s := New[string](2)
	s.Append("foo")
	s.InsertInIndex(0, "bas")
	s.InsertInIndex(1, "boo")
	s.RemoveFromIndex(1)
	fmt.Println(s.Pop())
	fmt.Println(s.All())
}

func main() {
	simpleArray()
	customSliceDriver()
}
