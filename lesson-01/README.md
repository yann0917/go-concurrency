# Mutex 如何解决资源并发访问问题

> 互斥锁 Mutex 就提供两个方法 Lock 和 Unlock：
>
> 进入临界区之前调用 Lock 方法，
>
> 退出临界区的时候调用 Unlock 方法

## Tools

* `go run -race xxx.go` 在编译（compile）、测试（test）或者运行（run）Go 代码的时候，加上 race 参数，就有可能发现并发问题.
* `go tool compile -race -S xxx.go` 查看编译后的汇编代码

## Attention

* 绝对不要把带race参数编译的程序部署到线上！
* 绝对不要把带race参数编译的程序部署到线上！
* 绝对不要把带race参数编译的程序部署到线上！

## Q&A

Q1. 如果 Mutex 已经被一个 goroutine 获取了锁，其它等待中的 goroutine 们只能一直等待。那么，等这个锁释放后，等待中的 goroutine 中哪一个会优先获取 Mutex 呢？

A1. [Mutex fairness](https://github.com/golang/go/blob/b94346e69bb01e1cd522ddfa9d09f41d9d4d3e98/src/sync/mutex.go#L42)
> 互斥锁有两种状态：正常状态和饥饿状态。
>
> 在正常状态下，所有等待锁的goroutine按照 `FIFO` 顺序等待。唤醒的goroutine不会直接拥有锁，而是会和新请求锁的goroutine竞争锁的拥有。新请求锁的goroutine具有优势：它正在CPU上执行，而且可能有好几个，所以刚刚唤醒的goroutine有很大可能在锁竞争中失败。在这种情况下，这个被唤醒的goroutine会加入到等待队列的前面。 如果一个等待的goroutine超过1ms没有获取锁，那么它将会把锁转变为饥饿模式。
>
> 在饥饿模式下，锁的所有权将从unlock的gorutine直接交给交给等待队列中的第一个。新来的goroutine将不会尝试去获得锁，即使锁看起来是unlock状态, 也不会去尝试自旋操作，而是放在等待队列的尾部。
>
> 如果一个等待的goroutine获取了锁，并且满足一以下其中的任何一个条件：(1)它是队列中的最后一个；(2)它等待的时候小于1ms。它会将锁的状态转换为正常状态。
>
> 正常状态有很好的性能表现，饥饿模式也是非常重要的，因为它能阻止尾部延迟的现象。

## Reference

* [package sync/mutex](https://golang.org/src/sync/mutex.go)
* [sync.mutex 源代码分析](https://colobu.com/2018/12/18/dive-into-sync-mutex/)

---
