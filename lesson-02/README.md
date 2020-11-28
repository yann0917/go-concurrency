# Mutex：庖丁解牛看实现

> CAS(compare-and-swap) 指令是实现互斥锁和同步原语的基础。CAS需要有3个操作数：内存地址V，旧的预期值A，即将要更新的目标值B。
>
> CAS指令执行时，当且仅当内存地址V的值与预期值A相等时，将内存地址V的值修改为B，否则就什么都不做。整个比较并替换的操作是一个原子操作。

## Mutex 的架构演进

### 0X001 初版

* 使用一个 flag 来表示锁是否被持有
* 问题：请求锁的 goroutine 会排队等待获取互斥锁。虽然这貌似很公平，但是从性能上来看，却不是最优的。因为如果我们能够把锁交给正在占用 CPU 时间片的 goroutine 的话，那就不需要做上下文的切换，在高并发的情况下，可能会有更好的性能。

```go
// 互斥锁的结构，包含两个字段
type Mutex struct {
    key int32 // 锁是否被持有的标识，0-锁未被持有，1-锁被持有，没有等待者，n-锁被持有，还有 n-1 个等待者
    sema int32 // 信号量专用，用以阻塞/唤醒goroutine
}
```

### 0X010 给新人机会

* 新的 goroutine 也尽可能地先获取到锁

```go

type Mutex struct {
    state int32
    sema  uint32
}


const (
    mutexLocked = 1 << iota // mutex is locked
    mutexWoken
    mutexWaiterShift = iota
)
```

state 是复合型字段：

1. mutexWaiters 表示阻塞等待此锁的 goroutine 数
2. mutexWoken 唤醒标记，代表是否有唤醒的 goroutine
3. mutexLocked 表示这个锁是否被持有

### 0X011 多给些机会

* 新来的和被唤醒的 goroutine 有更多的机会获取竞争锁
* 如果新来的 goroutine 或者是被唤醒的 goroutine 首次获取不到锁，它们就会通过自旋（spin，通过循环不断尝试，spin 的逻辑是在runtime 实现的）的方式，尝试检查锁是否被释放。在尝试一定的自旋次数后，再执行原来的逻辑。

### 0X100 解决饥饿

* 解决饥饿问题，不会让 goroutine 一直等待
* **Mutex 绝不容忍一个 goroutine 被落下，永远没有机会获取锁。** 不抛弃不放弃是它的宗旨，而且它也尽可能地让等待较长的 goroutine 更有机会获取到锁。

```go
type Mutex struct {
    state int32
    sema uint32
}

const (
    mutexLocked = 1 << iota // mutex is locked
    mutexWoken
    mutexStarving // 从state字段中分出一个饥饿标记
    mutexWaiterShift = iota
    starvationThresholdNs = 1e6
    )
```

state 是复合型字段：

1. mutexWaiters 表示阻塞等待此锁的 goroutine 数
2. mutexStarving 饥饿标记
3. mutexWoken 唤醒标记，代表是否有唤醒的 goroutine
4. mutexLocked 表示这个锁是否被持有

## Attention

1. **Unlock 方法可以被任意的 goroutine 调用释放锁，即使是没持有这个互斥锁的 goroutine，也可以进行这个操作。** 这是因为，Mutex 本身并没有包含持有这把锁的 goroutine 的信息，所以，Unlock 也不会对此进行检查。Mutex 的这个设计一直保持至今。
2. 一定要遵循『**谁申请，谁释放**』的原则，避免在 `foo()` 方法中申请加锁，在 `bar()` 方法中释放锁。

## Q&A

Q1.目前 Mutex 的 state 字段有几个意义，这几个意义分别是由哪些字段表示的？

A1.和第四阶段基本一致

Q2. 等待一个 Mutex 的 goroutine 数最大是多少？是否能满足现实的需求？

A2. 单从程序来看，可以支持 `1<<(32-3) -1` ，约 0.5 Billion 个，其中32为state的类型int32，3为waiter字段的shift。
考虑到实际 goroutine 初始化的空间为2K，0.5Billin*2K 达到了 1TB，单从内存空间来说已经要求极高了，当前的设计肯定可以满足了。

## References

* [CAS原理](https://www.jianshu.com/p/ab2c8fce878b)
* [互斥锁](https://golang.design/under-the-hood/zh-cn/part1basic/ch05sync/mutex/)

---
