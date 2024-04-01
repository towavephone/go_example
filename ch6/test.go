package main

import (
	"fmt"
	"go_example/ch6/geometry"
)

// go run ./ch6
func main() {
	perim := geometry.Path{{1, 1}, {5, 1}, {5, 4}, {1, 1}}
	fmt.Println(perim.Distance()) // "12", method of geometry.Path
}
