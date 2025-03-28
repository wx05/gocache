package main

import (
	"flag"
	"gocache/final"
	"log"
	"net/http"
)

// 注意这里，第三个参数类型 是 GetterFunc 类型，这里在后面的 g.getter.get 实现了这个接口，在g.getter.get 里面直接就是f()运行

// 启动缓存服务器:创建HTTPPool,添加节点信息,注册到cache中,启动HTTP服务(共3个端口,8001,8002,8003)
func startCacheServer(addr string, addrs []string, c *final.Config) {
	peers := final.NewHTTPPool(addr, c)
	peers.Set(c, addrs...)

	//数据预加载实现
	task := final.PreTask{}
	groups := task.Run(c, peers)
	//注册所有group的节点
	for _, g := range groups {
		g.RegisterPeers(peers)
	} //注册到一致性哈希的节点

	log.Println("gocache is running at", addr)
	log.Fatal(http.ListenAndServe(addr[7:], peers)) //从http:// 后面的localhost:8001 这里开始的
}

// 启动一个api服务(端口9999),与用户进行交互,用户通过访问api服务,获取缓存值，这里包括了从本地获取和从远程获取
func startAPIServer(apiAddr string, config *final.Config) {
	//后续修改成为restful api 形式
	http.Handle("/api", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			key := r.URL.Query().Get("key")
			groupName := r.URL.Query().Get("group")
			if len(groupName) == 0 {
				groupName = config.Server.DefalutGroup
			}
			group := final.GetGroup(groupName)
			view, err := group.Get(key) //
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
	var configPath string
	flag.IntVar(&port, "port", 8001, "gocache server port")
	flag.BoolVar(&api, "api", false, "start a api server")
	flag.StringVar(&configPath, "config", "./config.yaml", "config path")
	flag.Parse()
	//解析配置文件
	config, err := final.LoadConfig(configPath)
	if err != nil {
		log.Panicf("load config error, %s", err.Error())
		return
	}

	apiAddr := "http://localhost:9999" //对外提供服务的9999端口
	addrMap := map[int]string{         //缓存服务器的端口
		8001: "http://localhost:8001",
		8002: "http://localhost:8002",
		8003: "http://localhost:8003",
	}

	var addrs []string
	for _, v := range addrMap {
		addrs = append(addrs, v)
	}

	if api { //是否对外提供api接口服务
		go startAPIServer(apiAddr, config) //如果提供api接口对外服务,那么就起api接口
	}

	startCacheServer(addrMap[port], addrs, config) //8001端口是cache服务的端口
}
