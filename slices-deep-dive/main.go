package main

import "fmt"

func simpleArray() {
	var a [5]int
	a[0] = 1
	fmt.Println(a[0]) // Output: 1
	fmt.Println(a)    // Output: [1 0 0 0 0] , values are already initialized to 0 in
	// the capacity

}

func main() {
	simpleArray()
}
