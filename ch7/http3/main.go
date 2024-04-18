// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 194.

// Http3 is an e-commerce server that registers the /list and /price
// endpoints by calling (*http.ServeMux).Handle.
package main

import (
	"fmt"
	"log"
	"net/http"
)

type dollars float32

func (d dollars) String() string { return fmt.Sprintf("$%.2f", d) }

// !+main
// go run ./ch7/http2
// 访问 http://localhost:8002/list
// 访问 http://localhost:8002/price?item=shoes
func main() {
	db := database{"shoes": 50, "socks": 5}
	mux := http.NewServeMux()
	// db.list 的调用会援引一个接收者是 db 的 database.list 方法。
	// 所以 db.list 是一个实现了 handler 类似行为的函数，但是因为它没有方法（理解：该方法没有它自己的方法），所以它不满足 http.Handler 接口并且不能直接传给 mux.Handle
	// 语句 http.HandlerFunc(db.list) 是一个转换而非一个函数调用，因为 http.HandlerFunc 是一个类型
	mux.Handle("/list", http.HandlerFunc(db.list))
	mux.Handle("/price", http.HandlerFunc(db.price))
	log.Fatal(http.ListenAndServe("localhost:8002", mux))
}

type database map[string]dollars

func (db database) list(w http.ResponseWriter, req *http.Request) {
	for item, price := range db {
		fmt.Fprintf(w, "%s: %s\n", item, price)
	}
}

func (db database) price(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	price, ok := db[item]
	if !ok {
		w.WriteHeader(http.StatusNotFound) // 404
		fmt.Fprintf(w, "no such item: %q\n", item)
		return
	}
	fmt.Fprintf(w, "%s\n", price)
}

//!-main

/*
//!+handlerfunc
package http

type HandlerFunc func(w ResponseWriter, r *Request)

func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) {
	f(w, r)
}
//!-handlerfunc
*/
