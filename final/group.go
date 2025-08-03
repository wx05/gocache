package final

import (
	"fmt"
	"gocache/day6/singleflight"
	"log"
	"sync"
)

/*
	是

接收 key --> 检查是否被缓存 -----> 返回缓存值 ⑴

	|  否                         是
	|-----> 是否应当从远程节点获取 -----> 与远程节点交互 --> 返回缓存值 ⑵
	            |  否
	            |-----> 调用`回调函数`，获取值并添加到缓存 --> 返回缓存值 ⑶
*/
type Getter interface {
	Get(key string) ([]byte, error)
}

type GetterFunc func(key string) ([]byte, error)

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

type Group struct {
	name      string
	getter    Getter
	mainCache cache
	peer      peerPicker
	loader    *singleflight.Group
}

var (
	mu    sync.RWMutex
	group = make(map[string]*Group)
)

func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("getter func nil")
	}

	mu.Lock()
	defer mu.Unlock()

	g := &Group{
		name:      name,
		mainCache: cache{cacheBytes: cacheBytes},
		getter:    getter,
		loader:    &singleflight.Group{},
	}

	group[name] = g
	return g
}

func CreateGroup(c *Config, g GroupData) *Group {
	return NewGroup(g.Group, c.Server.MaxCacheBytes, GetterFunc(func(key string) ([]byte, error) {
		log.Println("[SlowDb] search key", key)
		if v, ok := Db[key]; ok {
			return []byte(v), nil
		}
		return nil, fmt.Errorf("%s not exists", key)
	}))
}

func GetGroup(name string) *Group {
	mu.RLock()
	g := group[name]
	mu.RUnlock()
	return g
}

func (g *Group) Get(key string) (ByteView, error) {

	if key == "" {
		return ByteView{}, fmt.Errorf("key can't be nil")
	}

	if v, ok := g.mainCache.get(key); ok == nil && v.Len() != 0 { //从本地的lru列表获取
		return v, ok
	}

	return g.load(key) //不存在则去加载
}

func (g *Group) RegisterPeers(peers peerPicker) {
	if g.peer != nil {
		panic("RegisterPeerPicker called more than once")
	}
	g.peer = peers
}

func (g *Group) load(key string) (value ByteView, err error) {

	//Do里面加锁，这样是的Do的第二个参数只被执行一次
	//第二个参数的作用是根据key求哈希值去节点获取数据
	viewi, err := g.loader.Do(key, func() (interface{}, error) {
		if g.peer != nil {
			if peer, ok := g.peer.PickPeer(key); ok { //选取节点
				if value, err := g.getFromPeer(peer, key); err == nil { //去节点获取值
					return value, err
				}
				log.Println("[GoCache] Failed to get from peer", err)
			}
		}
		return g.Locally(key) //如果在没有，那就去本地获取
	})

	if err == nil {
		return viewi.(ByteView), nil
	}

	return
}

func (g *Group) getFromPeer(peer PeerGetter, key string) (ByteView, error) {
	bytes, err := peer.Get(g.name, key)
	if err != nil {
		return ByteView{}, err
	}

	return ByteView{b: bytes}, err
}

func (g *Group) Locally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}

	value := ByteView{b: cloneBytes(bytes)}
	g.populateCache(key, value)
	return value, err
}

func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}
