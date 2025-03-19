package main

import (
	"fmt"
	"gocache/day3"
	"log"
	"net/http"
)

var db = map[string]string{
	"Tom": "20",
	"AAA": "90",
	"DDD": "100",
}

func main() {
	day3.NewGroup("scores", 2<<10, day3.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		},
	))

	addr := "localhost:9999"
	peers := day3.NewHTTPPool(addr)
	log.Println("geecache is running at", addr)
	log.Fatal(http.ListenAndServe(addr, peers))
}
