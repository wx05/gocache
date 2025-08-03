package day6

import (
	"fmt"
	"gocache/day6/consistenthash"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

func (p *HTTPPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format, v...))
}

type httpGetter struct {
	baseURL string
}

// 实现httpGetter接口
func (h *httpGetter) Get(group string, key string) ([]byte, error) {
	u := fmt.Sprintf(
		"%v%v/%v",
		h.baseURL,
		url.QueryEscape(group),
		url.QueryEscape(key),
	)

	res, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server resturn :%v", res.Status)
	}

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

var _ PeerGetter = (*httpGetter)(nil)

const (
	defaultBasePath = "/_gocache/"
	defaultReplicas = 50
)

type HTTPPool struct {
	self        string // self，用来记录自己的地址，包括主机名/IP 和端口
	basePath    string // basePath，作为节点间通讯地址的前缀，默认是 /_gocache/
	mu          sync.Mutex
	peers       *consistenthash.Map    //一致性哈希算法的 Map，用来根据具体的 key 选择节点
	httpGetters map[string]*httpGetter //映射远程节点与对应的 httpGetter
}

func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self:     self,
		basePath: defaultBasePath,
	}
}

/*
Set() 方法实例化了一致性哈希算法，并且添加了传入的节点。
并为每一个节点创建了一个 HTTP 客户端 httpGetter
*/
func (p *HTTPPool) Set(peers ...string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.peers = consistenthash.New(defaultReplicas, nil)
	p.peers.Add(peers...)
	p.httpGetters = make(map[string]*httpGetter, len(peers))
	for _, peer := range peers {
		p.httpGetters[peer] = &httpGetter{baseURL: peer + p.basePath}
	}
}

/*
PickerPeer() 包装了一致性哈希算法的 Get() 方法，根据具体的 key，选择节点，返回节点对应的 HTTP 客户端
*/
func (p *HTTPPool) PickPeer(key string) (PeerGetter, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if peer := p.peers.Get(key); peer != "" && peer != p.self {
		p.Log("Pick peer %s", peer)
		return p.httpGetters[peer], true
	}

	return nil, false
}

var _ peerPicker = (*HTTPPool)(nil) //鸭子类型（Duck Typing） 设计模式，只要一个类型实现了接口的所有方法，Go语言就认为该类型实现了该接口。也可以进行隐形类型转换。

// 缓存服务器对外的接口
/*
ServeHTTP 的实现逻辑是比较简单的，首先判断访问路径的前缀是否是 basePath，不是返回错误。
然后从路径中提取出 group 名称和 key，调用 group.Get(key) 获取缓存数据。
最后将缓存数据以 http.ResponseWriter 的形式返回。

这里实现 ServeHTTP ，进行主要的http请求过来后的逻辑
*/
func (p *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		panic("HTTPPool serving unexpected path:" + r.URL.Path)
	}

	p.Log("%s %s", r.Method, r.URL.Path)

	parts := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	groupName := parts[0]
	key := parts[1]

	group := GetGroup(groupName)
	if group == nil {
		http.Error(w, "no such group: "+groupName, http.StatusNotFound)
		return
	}

	view, err := group.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	_, err = w.Write(view.ByteSlice()) //注意这里只是将之前的clone了一次，而不是将view传过去,view只是做只读层
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
