package main

import (
	"fmt"
	"testing"
)

func TestTokenRecursiveMutex(t *testing.T) {
	r := &TokenRecursiveMutex{}
	a(r, 10086)
}

func a(l *TokenRecursiveMutex, token int64) {
	l.Lock(token)
	fmt.Printf("a---->")
	b(l, token)
	l.Unlock(token)
}
func b(l *TokenRecursiveMutex, token int64) {
	l.Lock(token)
	fmt.Printf("b---->")
	c(l, token)
	l.Unlock(token)
}
func c(l *TokenRecursiveMutex, token int64) {
	l.Lock(token)
	fmt.Printf("c\n")
	l.Unlock(token)
}
