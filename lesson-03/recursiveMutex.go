package main

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/petermattis/goid"
)

// RecursiveMutex 包装一个Mutex,实现可重入
type RecursiveMutex struct {
	sync.Mutex
	owner     int64 // 当前持有锁的goroutine id
	recursion int32 // 这个goroutine 重入的次数
}

// Lock lock
func (m *RecursiveMutex) Lock() {
	gID := goid.Get()
	// 如果当前持有锁的goroutine就是这次调用的goroutine,说明是重入
	if atomic.LoadInt64(&m.owner) == gID {
		m.recursion++
		return
	}
	m.Mutex.Lock()
	// 获得锁的goroutine第一次调用，记录下它的goroutine id,调用次数加1
	atomic.StoreInt64(&m.owner, gID)
	m.recursion = 1
}

// Unlock unlock
func (m *RecursiveMutex) Unlock() {
	gID := goid.Get()
	// 非持有锁的goroutine尝试释放锁，错误的使用
	if atomic.LoadInt64(&m.owner) != gID {
		panic(fmt.Sprintf("wrong the owner(%d):%d", m.owner, gID))
	}
	m.recursion--
	if m.recursion != 0 {
		return
	}
	// 此goroutine最后一次调用，需要释放锁
	atomic.StoreInt64(&m.owner, -1)
	m.Mutex.Unlock()
}

// TokenRecursiveMutex token 方式的递归锁
// Go 开发者不期望你利用 goroutine id 做一些不确定的东西，所以，他们没有暴露获取 goroutine id 的方法。
type TokenRecursiveMutex struct {
	sync.Mutex
	token     int64 // 当前持有锁的token
	recursion int32 // 这个goroutine 重入的次数
}

// Lock lock
func (m *TokenRecursiveMutex) Lock(token int64) {
	// 如果传入的token和持有锁的token一致，说明是递归调用
	if atomic.LoadInt64(&m.token) == token {
		m.recursion++
		return
	}
	m.Mutex.Lock()
	// 传入的token不一致，说明不是递归调用，抢到锁之后，记录这个 token
	atomic.StoreInt64(&m.token, token)
	m.recursion = 1
}

// Unlock unlock
func (m *TokenRecursiveMutex) Unlock(token int64) {
	// 非持有锁的goroutine尝试释放锁，错误的使用
	if atomic.LoadInt64(&m.token) != token {
		panic(fmt.Sprintf("wrong the owner(%d):%d", m.token, token))
	}
	m.recursion--
	if m.recursion != 0 {
		return
	}
	//  没有递归调用，需要释放锁
	atomic.StoreInt64(&m.token, 0)
	m.Mutex.Unlock()
}
