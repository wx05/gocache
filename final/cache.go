package final

import (
	"fmt"
	"gocache/day6/lru"
	"sync"
)

type cache struct {
	mu         sync.Mutex
	lru        *lru.Cache
	cacheBytes int64
}

func (c *cache) add(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()

	//延迟初始化
	if c.lru == nil {
		c.lru = lru.New(c.cacheBytes, nil)
	}

	c.lru.Add(key, value)
}

func (c *cache) get(key string) (ByteView, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	//延迟初始化
	if c.lru == nil {
		return ByteView{}, fmt.Errorf("lru isn't init")
	}

	if value, ok := c.lru.Get(key); ok {
		return value.(ByteView), nil //get时 修改类型为ByteView类型
	}

	return ByteView{}, nil
}
