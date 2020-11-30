package main

import (
	"fmt"
	"sync"
)

func foo(l sync.Locker) {
	fmt.Println("in foo")
	l.Lock()
	bar(l) // mutex 不可重入
	l.Unlock()
}

func bar(l sync.Locker) {
	l.Lock()
	fmt.Println("in bar")
	l.Unlock()
}

func reentrantLock() {
	l := &sync.Mutex{}
	foo(l)
}

func main() {
	reentrantLock()
}
