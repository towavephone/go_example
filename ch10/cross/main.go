// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 295.

// The cross command prints the values of GOOS and GOARCH for this target.
package main

import (
	"fmt"
	"runtime"
)

//!+
// go build ./ch10/cross
// ./cross
//
// GOARCH=386 go build ./ch10/cross
// ./cross

// 查看文档
// go get golang.org/x/tools/cmd/godoc
// go install golang.org/x/tools/cmd/godoc
// 配置环境变量一般为：$HOME/go/bin
// godoc -http :8000
func main() {
	fmt.Println(runtime.GOOS, runtime.GOARCH)
}

//!-