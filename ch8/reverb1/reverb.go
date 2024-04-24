// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 223.

// Reverb1 is a TCP server that simulates an echo.
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

// !+
func echo(c net.Conn, shout string, delay time.Duration) {
	fmt.Fprintln(c, "\t", strings.ToUpper(shout))
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", shout)
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", strings.ToLower(shout))
}

func handleConn(c net.Conn) {
	input := bufio.NewScanner(c)
	for input.Scan() {
		echo(c, input.Text(), 1*time.Second)
	}
	// NOTE: ignoring potential errors from input.Err()
	c.Close()
}

//!-
// go run ./ch8/reverb1
// 在另一个终端执行 go run ./ch8/netcat2，输入内容，回车，等待一段时间，会有回声
// 输出结果如下：
// Hello?
//     HELLO?
//     Hello?
//     hello?
// Is there anybody there?
//     IS THERE ANYBODY THERE?
// Yooo-hooo!
//     Is there anybody there?
//     is there anybody there?
//     YOOO-HOOO!
//     Yooo-hooo!
//     yooo-hooo!
// ^D
// 注意客户端的第三次 shout 在前一个 shout 处理完成之前一直没有被处理，这貌似看起来不是特别“现实”。真实世界里的回响应该是会由三次 shout 的回声组合而成的
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
