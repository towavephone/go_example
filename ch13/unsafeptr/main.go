// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 357.

// Package unsafeptr demonstrates basic use of unsafe.Pointer.
package main

import (
	"fmt"
	"unsafe"
)

// go run ./ch13/unsafeptr/main.go
func main() {
	//!+main
	var x struct {
		a bool
		b int16
		c []int
	}

	// equivalent to pb := &x.b
	pb := (*int16)(unsafe.Pointer(
		uintptr(unsafe.Pointer(&x)) + unsafe.Offsetof(x.b)))
	*pb = 42

	fmt.Println(x.b) // "42"
	//!-main
}

// 有时候垃圾回收器会移动一些变量以降低内存碎片等问题。这类垃圾回收器被称为移动 GC。
// 当一个变量被移动，所有的保存该变量旧地址的指针必须同时被更新为变量移动后的新地址。
// 从垃圾收集器的视角来看，一个 unsafe.Pointer 是一个指向变量的指针，因此当变量被移动时对应的指针也必须被更新；
// 但是 uintptr 类型的临时变量只是一个普通的数字，所以其值不应该被改变。上面错误的代码因为引入一个非指针的临时变量 tmp，导致垃圾收集器无法正确识别这个是一个指向变量 x 的指针。
// 当第二个语句执行时，变量 x 可能已经被转移，这时候临时变量 tmp 也就不再是现在的 &x.b 地址。第三个向之前无效地址空间的赋值语句将彻底摧毁整个程序
/*
//!+wrong
	// NOTE: subtly incorrect!
	tmp := uintptr(unsafe.Pointer(&x)) + unsafe.Offsetof(x.b)
	pb := (*int16)(unsafe.Pointer(tmp))
	*pb = 42
//!-wrong
*/
