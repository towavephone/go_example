// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 276.

// Package memo provides a concurrency-safe memoization a function of
// type Func.  Requests for different keys run concurrently.
// Concurrent requests for the same key result in duplicate work.
package memo

import "sync"

type Memo struct {
	f     Func
	mu    sync.Mutex // guards cache
	cache map[string]result
}

type Func func(string) (interface{}, error)

type result struct {
	value interface{}
	err   error
}

func New(f Func) *Memo {
	return &Memo{f: f, cache: make(map[string]result)}
}

// !+
// 减少了锁范围，使得并发性能提升，但有一些 URL 被获取了多次
// 导致重复执行，f 函数执行慢导致其中一个结果会被另一个结果覆盖
func (memo *Memo) Get(key string) (value interface{}, err error) {
	memo.mu.Lock()
	res, ok := memo.cache[key]
	memo.mu.Unlock()
	if !ok {
		res.value, res.err = memo.f(key)

		// Between the two critical sections, several goroutines
		// may race to compute f(key) and update the map.
		memo.mu.Lock()
		memo.cache[key] = res
		memo.mu.Unlock()
	}
	return res.value, res.err
}

//!-
