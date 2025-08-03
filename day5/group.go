package day5

import (
	"fmt"
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
	}

	group[name] = g
	return g
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

	if v, ok := g.mainCache.get(key); ok == nil {
		return v, ok
	}

	return g.load(key)
}

func (g *Group) RegisterPeers(peers peerPicker) {
	if g.peer != nil {
		panic("RegisterPeerPicker called more than once")
	}
	g.peer = peers
}

func (g *Group) load(key string) (value ByteView, err error) {
	if g.peer != nil {
		if peer, ok := g.peer.PickPeer(key); ok {
			if value, err := g.getFromPeer(peer, key); err == nil {
				return value, err
			}
			log.Println("[GeeCache] Failed to get from peer", err)
		}
	}
	return g.Locally(key)
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
