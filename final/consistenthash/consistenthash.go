package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

/*
一致性哈希的实现
1、数据结构 hash table 和 有序数组，hash table 用作哈希环， 有序数组的作用在于方便查找
2、一致性哈希的难点在于 进行 节点剔除，需要将之前的节点数据 顺时针分配到下个节点
*/

type Map struct {
	replica  int            //每个物理节点有多少个虚拟节点
	keys     []int          //hash key 存储位置，这里会排好序，方便查找 哈希值存储
	hashMap  map[int]string //哈希环 哈希值->物理节点key
	hashFunc hashFunc       //哈希函数
}

type hashFunc func(data []byte) uint32

// New
func New(replica int, hashF hashFunc) *Map {
	m := &Map{
		replica: replica,
		keys:    make([]int, 0),
		hashMap: make(map[int]string),
	}

	if hashF != nil {
		m.hashFunc = hashF
	} else {
		m.hashFunc = crc32.ChecksumIEEE
	}

	return m
}

// Add
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replica; i++ {
			//key 和下标，然后转byte，再转int类型
			hash := int(m.hashFunc([]byte(key + strconv.Itoa(i))))
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = key
		}
		sort.Ints(m.keys) //为什么需要排序,主要是因为方便查找
	}

}

// Get
func (m *Map) Get(key string) string {
	if len(key) == 0 {
		return ""
	}

	hash := int(m.hashFunc([]byte(key)))
	//查找第一个比他大的值,这里是二分查找
	index := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})

	//其实这里取余操作可以去掉,目的只是为了让结果在哈希环内
	return m.hashMap[(m.keys[(index % len(m.keys))])]
}
