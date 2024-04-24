// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 224.

// Reverb2 is a TCP server that simulates an echo.
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

func echo(c net.Conn, shout string, delay time.Duration) {
	fmt.Fprintln(c, "\t", strings.ToUpper(shout))
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", shout)
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", strings.ToLower(shout))
}

// !+
func handleConn(c net.Conn) {
	input := bufio.NewScanner(c)
	for input.Scan() {
		go echo(c, input.Text(), 1*time.Second)
	}
	// NOTE: ignoring potential errors from input.Err()
	c.Close()
}

// !-
// go run ./ch8/reverb2
// 在另一个终端执行 go run ./ch8/netcat2，输入内容，回车，等待一段时间，会有回声
// 输出结果如下：
// Is there anybody there?
//
//	IS THERE ANYBODY THERE?
//
// Yooo-hooo!
//
//	Is there anybody there?
//	YOOO-HOOO!
//	is there anybody there?
//	Yooo-hooo!
//	yooo-hooo!
//
// ^D
// 此时客户端的第三次 shout 在前一个 shout 处理完成之前就已经被处理了，这样看起来更像是真实世界里的回响，即并发了
func main() {
	l, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Print(err) // e.g., connection aborted
			continue
		}
		go handleConn(conn)
	}
}
