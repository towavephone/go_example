// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 227.

// Netcat is a simple read/write client for TCP servers.
package main

import (
	"io"
	"log"
	"net"
	"os"
)

// !+
// 当用户关闭了标准输入，主 goroutine 中的 mustCopy 函数调用将返回，然后调用 conn.Close() 关闭读和写方向的网络连接。
// 关闭网络连接中的写方向的连接将导致 server 程序收到一个文件（end-of-file）结束的信号。
// 关闭网络连接中读方向的连接将导致后台 goroutine 的 io.Copy 函数调用返回一个“read from closed connection”（“从关闭的连接读”）类似的错误，因此我们临时移除了错误日志语句；
// 在练习 8.3 将会提供一个更好的解决方案。（需要注意的是 go 语句调用了一个函数字面量，这是 Go 语言中启动 goroutine 常用的形式。）
// 在后台 goroutine 返回之前，它先打印一个日志信息，然后向 done 对应的 channel 发送一个值。主 goroutine 在退出前先等待从 done 对应的 channel 接收一个值。
// 因此，总是可以在程序退出前正确输出 “done” 消息。
func main() {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	done := make(chan struct{})
	go func() {
		io.Copy(os.Stdout, conn) // NOTE: ignoring errors
		log.Println("done")
		done <- struct{}{} // signal the main goroutine
	}()
	mustCopy(conn, os.Stdin)
	conn.Close()
	<-done // wait for background goroutine to finish
}

//!-

func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}
