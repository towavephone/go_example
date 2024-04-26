// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 263.

// Package bank provides a concurrency-safe single-account bank.
package bank

//!+
import "sync"

var (
	mu      sync.Mutex // guards balance
	balance int
)

func Deposit(amount int) {
	mu.Lock()
	balance = balance + amount
	mu.Unlock()
}

func Balance() int {
	mu.Lock()
	b := balance
	mu.Unlock()
	return b
}

// 允许多个只读操作并行执行，但写操作会完全互斥
// 这种锁叫作“多读单写”锁（multiple readers, single writer lock）
// var mu sync.RWMutex
// var balance int
// func Balance() int {
//     mu.RLock() // readers lock
//     defer mu.RUnlock()
//     return balance
// }

//!-
// 保证取款操作的原子性，可以使用下面的代码：
// func Withdraw(amount int) bool {
//     mu.Lock()
//     defer mu.Unlock()
//     deposit(-amount)
//     if balance < 0 {
//         deposit(amount)
//         return false // insufficient funds
//     }
//     return true
// }

// func Deposit(amount int) {
//     mu.Lock()
//     defer mu.Unlock()
//     deposit(amount)
// }

// func Balance() int {
//     mu.Lock()
//     defer mu.Unlock()
//     return balance
// }

// // This function requires that the lock be held.
// func deposit(amount int) { balance += amount }

// 下面的　go　代码会输出什么？
// var x, y int
// go func() {
//     x = 1 // A1
//     fmt.Print("y:", y, " ") // A2
// }()
// go func() {
//     y = 1                   // B1
//     fmt.Print("x:", x, " ") // B2
// }()

// y:0 x:1
// x:0 y:1
// x:1 y:1
// y:1 x:1
// x:0 y:0
// y:0 x:0
// 其中最后两个输出的原因是：
// 多核cpu中，并发运行时，在编译器编译后 x, y 变量有可能是在两个独立的 CPU 上都有副本的，并且此时是被初始化为 0 的，
// 由于编译器认为 A1, A2 这两条语句和 B1, B2 这两条语句的顺序不影响结果，就有可能调换两者的次序。从而打印出 x:0 y:0，y:0 x:0

// 这里引入了多读单写锁保证读取性能，解决了惰性初始化问题
// func loadIcons() {
//     icons = make(map[string]image.Image)
//     icons["spades.png"] = loadIcon("spades.png")
//     icons["hearts.png"] = loadIcon("hearts.png")
//     icons["diamonds.png"] = loadIcon("diamonds.png")
//     icons["clubs.png"] = loadIcon("clubs.png")
// }
//
// var mu sync.RWMutex // guards icons
// var icons map[string]image.Image
// // Concurrency-safe.
// func Icon(name string) image.Image {
//     mu.RLock()
//     if icons != nil {
//         icon := icons[name]
//         mu.RUnlock()
//         return icon
//     }
//     mu.RUnlock()

//     // acquire an exclusive lock
//     mu.Lock()
//     if icons == nil { // NOTE: must recheck for nil
//         loadIcons()
//     }
//     icon := icons[name]
//     mu.Unlock()
//     return icon
// }

// 当然有更好的方法，可以使用 sync.Once 来实现，等价于上面代码
// var loadIconsOnce sync.Once
// var icons map[string]image.Image
// // Concurrency-safe.
// func Icon(name string) image.Image {
//     loadIconsOnce.Do(loadIcons)
//     return icons[name]
// }
