package lru

import (
	"reflect"
	"testing"
)

// 测试空值是否可以放进去
func TestLRU_Empty(t *testing.T) {
	cache := New(100, nil)
	cache.Add("", "v") // 空key应被忽略
	cache.Add("k", "") // 空value应被忽略
	if cache.Get("k") != "" {
		t.Errorf("Should ignore empty value")
	}
}

/*
1、测试add是否ok
2、测试add之后是否可以get出来
*/
func TestCache_Add(t *testing.T) {
	lru := New(10, nil)
	key := "aaaaa"
	value := "bbbbb"
	lru.Add(key, value)

	res := lru.Get(key)
	if res != value {
		t.Errorf("Add error")
	}
}

/*
1、测试超过最大byte时，是否会淘汰最远没有被使用的
2、回调函数是否会被正确运行
*/
func TestCache_onEvicted(t *testing.T) {
	keys := make([]string, 0)
	onEvicted := func(key string, value string) {
		keys = append(keys, key)
	}
	lru := New(10, onEvicted)
	lru.Add("k1", "v1")
	lru.Add("k2", "v1")
	lru.Add("k3", "v1")
	lru.Add("k4", "v1")

	//expect := []string{"k2", "k1"}
	expect := []string{"k1", "k2"} //这两行确保是为了测试 k1先被置换出来，然后是k2

	if !reflect.DeepEqual(expect, keys) {
		t.Fatalf("Call OnEvicted failed, expect keys equals to %s", expect)
	}
}
