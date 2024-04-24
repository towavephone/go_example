// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 228.

// Pipeline1 demonstrates an infinite 3-stage pipeline.
package main

import "fmt"

//!+
// go run ./ch8/pipeline1
func main() {
	naturals := make(chan int)
	squares := make(chan int)

	// Counter
	go func() {
		for x := 0; ; x++ {
			naturals <- x
		}
	}()

	// Squarer
	go func() {
		for {
			x := <-naturals
			// 当 naturals 对应的 channel 被关闭并没有值可接收时跳出循环，并且也关闭 squares 对应的 channel
			// 等价于使用 range 循环
			// x, ok := <-naturals
			// if !ok {
			// 	break // channel was closed and drained
			// }
			squares <- x * x
		}
		// close(squares)
	}()

	// Printer (in main goroutine)
	for {
		fmt.Println(<-squares)
	}
}

//!-