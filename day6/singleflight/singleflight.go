package singleflight

import "sync"

/*
call 代表正在进行中，或已经结束的请求。使用 sync.WaitGroup 锁避免重入
*/
type call struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}

/*
Group 是 singleflight 的主数据结构，管理不同 key 的请求(call)
*/
type Group struct {
	mu sync.Mutex
	m  map[string]*call
}

/*
针对相同的 key，无论 Do 被调用多少次，函数 fn 都只会被调用一次，等待 fn 调用结束了，返回返回值或错误
*/
func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {

	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}

	//如果当前的key正在运行,则wait并且返回结果
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		c.wg.Wait() //等到所有的都被释放了
		return c.val, c.err
	}

	c := new(call)
	c.wg.Add(1) //请求前加锁,表明已经有key在进行处理
	g.m[key] = c
	g.mu.Unlock()

	c.val, c.err = fn() //调用函数，发起请求
	c.wg.Done()         //解锁

	g.mu.Lock()
	delete(g.m, key) //更新 g.m
	g.mu.Unlock()

	return c.val, c.err
}
