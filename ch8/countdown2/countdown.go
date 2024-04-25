// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 244.

// Countdown implements the countdown for a rocket launch.
package main

import (
	"fmt"
	"os"
	"time"
)

// !+
// go run ./ch8/countdown2
func main() {
	// ...create abort channel...

	//!-

	//!+abort
	abort := make(chan struct{})
	go func() {
		os.Stdin.Read(make([]byte, 1)) // read a single byte
		abort <- struct{}{}
	}()
	//!-abort

	//!+
	fmt.Println("Commencing countdown... Press return to abort.")
	// select 会等待 case 中有能够执行的 case 时去执行。当条件满足时，select 才会去通信并执行 case 之后的语句；这时候其它通信是不会执行的。
	// 一个没有任何 case 的 select 语句写作 select{}，会永远地等待下去。
	select {
	case <-time.After(10 * time.Second):
		// Do nothing.

	case <-abort:
		fmt.Println("Launch aborted!")
		return
	}
	launch()
}

//!-

func launch() {
	fmt.Println("Lift off!")
}
