# day1

源网址：https://geektutu.com/post/geecache-day1.html

## lru 实现
思路：
1、使用双向链表 和 hash table 来实现，其中，hash table 主要是 查找效率为O(1), hash table 则为了方便查找 lru链表里面的数据节点 

2、双向链表里面存储的是entry结构体，里面存储具体的数据，hash table key是存储的key，value是指向 hash table node 节点的指针，双向链表的作用是为了挪动节点，进行最远未使用的节点的淘汰

## 差异点:
1、我在实现时，考虑到线程安全，故而，在Cache结构体里面加了一个 lock sync.Mutex 的成员，来加锁

## 待优化点
在Get时由于对于 list修改，因为加了写锁，这会导致效率变慢，优化点方案共有三:

1、在读取时加读锁，不使用defer, 在更改时加写锁，不使用defer

2、将锁的粒度变小，分段锁，而不是对于整个cache进行加锁。使用hash 函数(fnv32)对于缓存分段，然后使用和方案1一样的加锁方式

3、原子操作和 CAS，比较复杂，需要使用 sync/atomic 来实现，由于 container/list 不是线程安全的，所以需要自己使用sync/atomic包来自己实现list