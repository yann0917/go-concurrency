package main

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {

	os.Exit(m.Run())
}

func TestNoLockCounter(t *testing.T) {
	NoLockCounter()
}

func TestLockCounter(t *testing.T) {
	LockCounter()
}

func TestCounter(t *testing.T) {
	t.Run("GetCounter", func(t *testing.T) {
		count := GetCounter()
		t.Log(count)
	})
}
