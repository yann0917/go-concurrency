package main

import (
	"fmt"
	"sync"
)

// NoLockCounter 未加互斥锁，导致数据竞争
func NoLockCounter() {
	var count = 0
	var wg sync.WaitGroup
	wg.Add(10)

	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			for i := 0; i < 10000; i++ {
				count++
			}
		}()
	}

	wg.Wait()
	fmt.Printf("noLock: %v", count)
	fmt.Println()
}

// LockCounter 加互斥锁
func LockCounter() {
	var count = 0
	var mu sync.Mutex
	var wg sync.WaitGroup
	wg.Add(10)

	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			for i := 0; i < 10000; i++ {
				mu.Lock()
				count++
				mu.Unlock()
			}
		}()
	}

	wg.Wait()
	fmt.Printf("lock: %v", count)
	fmt.Println()
}

// Counter mutex 嵌入结构体，线程安全，不对外暴露加锁逻辑
type Counter struct {
	mu    sync.Mutex
	count uint64
}

// Incr count++
func (c *Counter) Incr() {
	c.mu.Lock()
	c.count++
	c.mu.Unlock()
}

// Count get count
func (c *Counter) Count() uint64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.count
}

// GetCounter mutex 嵌入结构体，隐藏互斥锁逻辑
func GetCounter() uint64 {
	var counter Counter
	var wg sync.WaitGroup
	wg.Add(10)

	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < 10000; j++ {
				counter.Incr()
			}
		}()
	}
	wg.Wait()
	fmt.Print(counter.Count())
	return counter.Count()
}
func main() {
	NoLockCounter()
	LockCounter()
	GetCounter()
}
