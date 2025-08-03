package lru

import (
	"container/list"
)

/*
主要思路：
一个哈希hash，主要做数据查找使用
一个双向链表，维护元素的访问顺序（链表头为最近使用，尾部为最久未使用）。
*/

type Cache struct {
	maxBytes  int64                         //最大容量
	nBytes    int64                         //当前容量
	ll        *list.List                    //lru List，最新的节点总是在最前面,
	cache     map[string]*list.Element      //数据存储/最近数据总是在list头
	onEvicted func(key string, value Value) //对于淘汰的数据处理时的回调函数
}

// ll 里面的成员结构
type entry struct {
	key   string
	value Value
}

type Value interface {
	Len() int
}

// New 创建lru列表
func New(maxBytes int64, onEvicted func(key string, value Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		onEvicted: onEvicted,
	}
}

// Add key vale set
func (c *Cache) Add(key string, value Value) {
	//新增对于空值的处理
	if len(key) == 0 {
		return
	}

	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(c.cache[key])
		kv := ele.Value.(*entry)
		kv.value = value
		c.cache[key] = c.ll.Front() //因为前面已经push到Front了，只需要更新索引即可。

		diff := int64(ele.Value.(*entry).value.Len() - value.Len()) //计算出当前的value的差异
		c.nBytes += diff                                            //更新当前的容量

	} else {
		//新的节点放入
		node := &entry{
			key:   key,
			value: value,
		}
		c.cache[key] = c.ll.PushFront(node)
		c.nBytes += int64(len(key) + value.Len())

	}

	//循环删除,如果不够就一直删下去（自动缩容）
	for c.maxBytes > 0 && c.nBytes > c.maxBytes {
		c.RemoveOldData()
	}
}

func (c *Cache) RemoveOldData() {

	ele := c.ll.Back()
	if ele != nil {
		key := ele.Value.(*entry).key
		value := ele.Value.(*entry).value

		c.nBytes -= int64(len(key) + value.Len())
		if c.nBytes <= 0 {
			c.nBytes = 0
		}

		delete(c.cache, key)
		c.ll.Remove(ele)

		if c.onEvicted != nil {
			c.onEvicted(key, value)
		}
	}
}

func (c *Cache) Get(key string) (value Value, ok bool) {

	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(c.cache[key])
		kv := ele.Value.(*entry)
		return kv.value, true
	}

	return
}

func (c *Cache) Len() int {
	return c.ll.Len()
}
