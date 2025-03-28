package main

import (
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

type RegisterCenter struct {
	mu    sync.RWMutex
	nodes map[string]time.Time
}

// 初始化注册中心
func NewRegisterCenter() *RegisterCenter {
	return &RegisterCenter{
		nodes: make(map[string]time.Time),
	}
}

// 注册中心逻辑
func (r *RegisterCenter) handleServe(writer http.ResponseWriter, request *http.Request) {
	addr, err := io.ReadAll(io.LimitReader(request.Body, 1024)) //最多也就1KB的数据
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		log.Print("read request body error")
		return
	}

	if len(addr) == 0 {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Printf(" receive request addr : %s", addr)

	r.mu.Lock()
	r.nodes[string(addr)] = time.Now()
	r.mu.Unlock()

	writer.WriteHeader(http.StatusOK)
}

func (r *RegisterCenter) checkNodesHealth() {
	for {
		r.mu.RLock()
		for addr, kickTime := range r.nodes {
			log.Println(time.Since(kickTime) > 5*time.Second, addr)
			if time.Since(kickTime) > 5*time.Second {
				r.mu.Lock()
				delete(r.nodes, addr)
				r.mu.Unlock()
				log.Printf("[ERROR] nodes offline : %s, last kick time: %s", addr, kickTime)
			}
		}
		r.mu.RUnlock()
		time.Sleep(3 * time.Second) //sleep 3秒之后再运行
	}
}

func main() {
	rc := NewRegisterCenter()
	http.HandleFunc("/register", rc.handleServe)
	log.Print(" Running at port:8000...")
	go rc.checkNodesHealth()
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatalf("start server error: %s", err.Error())
		return
	}

}
