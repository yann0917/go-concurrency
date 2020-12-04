# RWMutex：读写锁的实现原理及避坑指南

> 标准库中的 RWMutex 是一个 reader/writer 互斥锁。RWMutex 在某一时刻只能由任意数量的 reader 持有，或者是只被单个的 writer 持有。

## RWMutex 的实现原理

读写锁的设计和实现分成三类：

1. Read-preferring：读优先的设计可以提供很高的并发性，但是，在竞争激烈的情况下可能会导致写饥饿。
2. Write-preferring：写优先的设计意味着，如果已经有一个 writer 在等待请求锁的话，它会阻止新来的请求锁的 reader 获取到锁，所以优先保障 writer。
3. 不指定优先级：这种设计比较简单，不区分 reader 和 writer 优先级，某些场景下这种不指定优先级的设计反而更有效，因为第一类优先级会导致写饥饿，第二类优先级可能会导致读饥饿，这种不指定优先级的访问不再区分读写，大家都是同一个优先级，解决了饥饿的问题。

* Go 标准库中的 RWMutex 是基于 `Mutex` 实现的。
* Go 标准库中的 RWMutex 设计是 `Write-preferring` 方案。一个正在阻塞的 Lock 调用会排除新的 reader 请求到锁。

```go
type RWMutex struct {
  w           Mutex   // 互斥锁解决多个writer的竞争
  writerSem   uint32  // writer阻塞信号量
  readerSem   uint32  // reader阻塞信号量
  readerCount int32   // 当前reader的数量（以及是否有 writer 竞争锁）
  readerWait  int32   // writer等待完成的reader的数量
}

const rwmutexMaxReaders = 1 << 30 // 最大的 reader 数量
```

## RWMutex 的 3 个踩坑点

> 在使用读写锁的时候，一定要注意，不遗漏不多余。
>
> `Lock` 和 `RLock` 多余的调用会导致锁没有被释放，可能会出现死锁，而 `Unlock` 和 `RUnlock` 多余的调用会导致 panic。

### 不可复制

* 一旦读写锁被使用，它的字段就会记录它当前的一些状态。这个时候你去复制这把锁，就会把它的状态也给复制过来。但是，原来的锁在释放的时候，并不会修改你复制出来的这个读写锁，这就会导致复制出来的读写锁的状态不对，可能永远无法释放锁。

### 重入导致死锁

* 因为读写锁内部基于互斥锁实现对 writer 的并发访问，而互斥锁本身是有重入问题的，所以，writer 重入调用 Lock 的时候，就会出现死锁的现象。
* 有活跃 reader 的时候，writer 会等待，如果我们在 reader 的读操作时调用 writer 的写操作（它会调用 Lock 方法），那么，这个 reader 和 writer 就会形成互相依赖的死锁状态。
* writer 依赖活跃的 reader -> 活跃的 reader 依赖新来的 reader -> 新来的 reader 依赖 writer。

### 释放未加锁的 RWMutex

* Lock 和 Unlock 的调用总是成对出现的，RLock 和 RUnlock 的调用也必须成对出现。

## Attention

如果你能意识到你要解决的是一个 readers-writers 问题，那么你就可以毫不犹豫地选择 RWMutex，不用考虑其它选择。
如果你能意识到你要解决的是一个 readers-writers 问题，那么你就可以毫不犹豫地选择 RWMutex，不用考虑其它选择。
如果你能意识到你要解决的是一个 readers-writers 问题，那么你就可以毫不犹豫地选择 RWMutex，不用考虑其它选择。

## Q&A

Q1.

A1.

## References

* [readers-writers 问题](https://en.wikipedia.org/wiki/Readers%E2%80%93writers_problem)

---
