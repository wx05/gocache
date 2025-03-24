package main

import (
	"flag"
	"fmt"
	"gocache/day5"
	"log"
	"net/http"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func createGroup() *day5.Group {
	return day5.NewGroup("scores", 2<<10, day5.GetterFunc(func(key string) ([]byte, error) {
		log.Println("[SlowDb] search key", key)
		if v, ok := db[key]; ok {
			return []byte(v), nil
		}
		return nil, fmt.Errorf("%s not exists", key)
	}))
}

// 启动缓存服务器:创建HTTPPool,添加节点信息,注册到gocache中,启动HTTP服务(共3个端口,8001,8002,8003)
func startCacheServer(addr string, addrs []string, g *day5.Group) {
	peers := day5.NewHTTPPool(addr)
	peers.Set(addrs...)
	g.RegisterPeers(peers) //注册到一致性哈希的节点
	log.Println("gocache is running at", addr)
	log.Fatal(http.ListenAndServe(addr[7:], peers))
}

// 启动一个api服务(端口9999),与用户进行交互,用户通过访问api服务,获取缓存值，这里包括了从本地获取和从远程获取
func startAPIServer(apiAddr string, gee *day5.Group) {
	http.Handle("/api", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			key := r.URL.Query().Get("key")
			view, err := gee.Get(key) //从本地缓存中获取数据
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Write(view.ByteSlice())

		}))
	log.Println("fontend server is running at", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr[7:], nil))
}

func main() {
	var port int
	var api bool
	flag.IntVar(&port, "port", 8001, "gocache server port")
	flag.BoolVar(&api, "api", false, "start a api server")
	flag.Parse()

	apiAddr := "http://localhost:9999" //对外提供服务的9999端口
	addrMap := map[int]string{
		8001: "http://localhost:8001",
		8002: "http://localhost:8002",
		8003: "http://localhost:8003",
	}

	var addrs []string
	for _, v := range addrMap {
		addrs = append(addrs, v)
	}

	goCache := createGroup() //创建数据组
	if api {
		go startAPIServer(apiAddr, goCache) //如果提供api接口对外服务,那么就起api接口
	}

	startCacheServer(addrMap[port], []string(addrs), goCache) //8001端口是cache服务的端口
}
