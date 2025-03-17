package day2

import (
	"fmt"
	"log"
	"testing"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func TestGet(t *testing.T) {
	loadCounts := make(map[string]int, len(db))

	gee := NewGroup("test", 2<<10, GetterFunc(func(key string) ([]byte, error) {
		log.Println("[SlowDB] search key", key)
		if v, ok := db[key]; ok {
			if _, ok := loadCounts[key]; !ok {
				loadCounts[key] = 0
			}
			loadCounts[key]++
			return []byte(v), nil
		}
		return nil, fmt.Errorf("%s not exists", key)
	}))

	for k, v := range db {
		if view, err := gee.Get(k); err != nil || view.String() != v {
			//t.Fatalf("failed to get value:key:%s,value:%s", k, v)
		}

		if _, err := gee.Get(k); err != nil || loadCounts[k] > 1 {
			t.Fatalf("cache %s miss", k)
		}
		//cache hit
	}
	if view, err := gee.Get("aaa"); err != nil {
		t.Fatalf("the value of unknow should be empty, but %s got", view)
	}
}
