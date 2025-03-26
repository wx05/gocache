# day6
https://geektutu.com/post/geecache-day6.html

## 思路
防止缓存击穿

使用sync.WaitGroup 和 sync.Mutex 来进行加锁,保证只有一个进程在处理 

主要流程：
1、main.go 根据参数来判断是否开放api服务端口，默认开启cache server
2、请求api服务，则直接去本地缓存读取数据，本地如果不存在，则load并返回，这里没有节点的注册，仅在本地读取
3、请求cache server服务，将所有的节点都注册到一致性哈希环里面，并且每一个设置默认的虚拟节点数，在 ServeHTTP里面进行 接受请求时的处理
4、在 ServeHTTP 里面，获取数据分组，根据分组去获取数据，实现在  group.go 文件里面 的Get 函数里面，现在本地节点读取，如果不存在，则根据哈希算出来需要去哪个节点去获取数据，http.get里面去获取数据,然后返回


需要注意的点在于：
1、startCacheServer 接口返回的是一个 *HTTPPool类型，但是下面的 RegisterPeers 的参数是 peerPicker类型，一个interface类型，这个主要是因为 *HTTPPool实现了 peerPicker的所有类型，也叫鸭子类型，所以，这里可以直接传参进去，请注意这一行 var _ peerPicker = (*HTTPPool)(nil)  这里会进行编译器校验。传参进去也会进行隐式转化。