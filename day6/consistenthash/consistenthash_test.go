package consistenthash

import (
	"strconv"
	"testing"
)

func TestHashing(t *testing.T) {
	hashF := func(key []byte) uint32 {
		i, _ := strconv.Atoi(string(key))
		return uint32(i)
	}
	hash := New(3, hashF)

	//这里是将数字转为下标，然后每一个节点有三个子节点
	hash.Add("6", "4", "2")

	testCases := map[string]string{
		"2":  "2",
		"11": "2",
		"23": "4",
		"27": "2",
	}

	//testCases的map去Get，如果Get出来的子节点不一样，报错
	for k, v := range testCases {
		if hash.Get(k) != v {
			t.Errorf("Asking for %s, should have yielded %s", k, v)
		}
	}

	hash.Add("8")

	testCases["27"] = "8"

	for k, v := range testCases {
		if hash.Get(k) != v {
			t.Errorf("Asking for %s, should have yielded %s", k, v)
		}
	}
}
