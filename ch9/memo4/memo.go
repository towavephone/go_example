// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 276.

// Package memo provides a concurrency-safe memoization a function of
// a function.  Requests for different keys proceed in parallel.
// Concurrent requests for the same key block until the first completes.
// This implementation uses a Mutex.
package memo

import "sync"

type Memo struct {
	f     Func
	mu    sync.Mutex // guards cache
	cache map[string]*entry
}

// Func is the type of the function to memoize.
type Func func(string) (interface{}, error)

type result struct {
	value interface{}
	err   error
}

// !+
type entry struct {
	res   result
	ready chan struct{} // closed when res is ready
}

func New(f Func) *Memo {
	return &Memo{f: f, cache: make(map[string]*entry)}
}

func (memo *Memo) Get(key string) (value interface{}, err error) {
	memo.mu.Lock()
	e := memo.cache[key]
	if e == nil {
		// This is the first request for this key.
		// This goroutine becomes responsible for computing
		// the value and broadcasting the ready condition.
		// 插入一个未准备好的条目
		e = &entry{ready: make(chan struct{})}
		memo.cache[key] = e
		memo.mu.Unlock()

		e.res.value, e.res.err = memo.f(key)

		// 告诉这个条目准备好了
		close(e.ready) // broadcast ready condition
	} else {
		// This is a repeat request for this key.
		memo.mu.Unlock()

		<-e.ready // wait for ready condition
	}
	// 条目中的 e.res.value 和 e.res.err 变量是在多个 goroutine 之间共享的。
	// 创建条目的 goroutine 同时也会设置条目的值，其它 goroutine 在收到 ready 的广播消息之后立刻会去读取条目的值。
	// 尽管会被多个 goroutine 同时访问，但却并不需要互斥锁。
	// ready channel 的关闭一定会发生在其它 goroutine 接收到广播事件之前，因此第一个 goroutine 对这些变量的写操作是一定发生在这些读操作之前的。不会发生数据竞争。
	return e.res.value, e.res.err
}

//!-
