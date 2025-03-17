
# day2

https://geektutu.com/post/geecache-day2.html


## 思路
1、sync.Mutex 使用，实现并发控制， 和我day1想到的一样，需要使用锁来进行并发控制
只是在想锁的控制应该在lru的使用层 还是在lru层，这是个问题


2、group 结构体的实现,主要实现不同命名空间的，类似于 redis里面 数据库编号

