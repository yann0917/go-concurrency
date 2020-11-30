# Mutex：4种易错场景大盘点

## 常见错误

### Lock/Unlock 不是成对出现

1. 常见的缺少 `Unlock()` 的场景
   1. 代码中有太多的 if-else 分支，可能在某个分支中漏写了 Unlock
   2. 在重构的时候把 Unolock 删除了
   3. Unlock 误写成 Lock

[参考代码](../lesson-01/counter.go)

### Copy 已使用的 Mutex

> Package sync 的同步原语在使用后是不能复制的
>
> 原因是 Mutex 是一个有状态的对象，它的 state 字段记录这个锁的状态。如果你要复制一个已经加锁的 Mutex 给一个新的变量，那么新的刚初始化的变量居然被加锁了，这显然不符合你的期望，因为你期望的是一个零值的 Mutex。关键是在并发环境下，你根本不知道要复制的 Mutex 状态是什么，因为要复制的 Mutex 是由其它 goroutine 并发访问的，状态可能总是在变化。

### 重入(递归锁)

> 当一个线程获取锁时，如果没有其它线程拥有这个锁，那么这个线程就成功获取到这个锁。之后，如果其它线程再请求这个锁，就会处于阻塞等待的状态。但是，如果拥有这把锁的线程再请求这把锁的话，不会阻塞，而是成功返回，所以叫可重入锁（有时候也叫做递归锁）。
>
> **Mutex 不是可重入的锁**。

[Mutex error case](reentrantLock.go)

如何实现一个可重入锁？

方案一：goroutine id

* 简单方式，就是通过 `runtime.Stack()` 方法获取栈帧信息，栈帧信息里包含 `goroutine id
* hacker 方式 [case](recursiveMutex.go)

方案二：token

* 调用者自己提供一个 token，获取锁的时候把这个 token 传入，释放锁的时候也需要把这个 token 传入。通过用户传入的 token 替换方案一中 goroutine id [case](recursiveMutex_test.go)

## 死锁

> 两个或两个以上的进程（或线程，goroutine）在执行过程中，因争夺共享资源而处于一种互相等待的状态，如果没有外部干涉，它们都将无法推进下去，此时，我们称系统处于死锁状态或系统产生了死锁。

如果在一个系统中以下四个条件同时成立，那么就能引起死锁：

1. 互斥： 至少一个资源是被排他性独享的，其他线程必须处于等待状态，直到资源被释放。
2. 持有和等待：goroutine 持有一个资源，并且还在请求其它 goroutine 持有的资源，也就是咱们常说的“吃着碗里，看着锅里”的意思。
3. 不可剥夺：资源只能由持有它的 goroutine 来释放。
4. 环路等待：一般来说，存在一组等待进程，P={P1，P2，…，PN}，P1 等待 P2 持有的资源，P2 等待 P3 持有的资源，依此类推，最后是 PN 等待 P1 持有的资源，这就形成了一个环路等待的死结。

## Tools

* [go-deadlock](https://github.com/sasha-s/go-deadlock) Online deadlock detection in go

* [go-tools](https://github.com/dominikh/go-tools) Staticcheck - The advanced Go linter
* [goid](https://github.com/petermattis/goid) Programatically retrieve the current goroutine's ID.

## Attention

* 保证 Lock/Unlock 成对出现，尽可能采用 defer mutex.Unlock 的方式，把它们成对、紧凑地写在一起。
* 保证 Lock/Unlock 成对出现，尽可能采用 defer mutex.Unlock 的方式，把它们成对、紧凑地写在一起。
* 保证 Lock/Unlock 成对出现，尽可能采用 defer mutex.Unlock 的方式，把它们成对、紧凑地写在一起。

## Q&A

Q1.查找知名的数据库系统 TiDB 的 issue，看看有没有 Mutex 相关的 issue，看看它们都是哪些相关的 Bug

A1. [TiDB issue deadlock](https://github.com/pingcap/tidb/issues?q=deadlock)

## References

* [死锁-哲学家就餐问题](https://zh.wikipedia.org/wiki/%E5%93%B2%E5%AD%A6%E5%AE%B6%E5%B0%B1%E9%A4%90%E9%97%AE%E9%A2%98)

---
