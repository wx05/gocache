
# day3
https://geektutu.com/post/geecache-day3.html

## 思路
开放一个http端口，可以与外部进行交互

## 坑
第三天这里遇到了go包管理的一些问题，在此现将一些概念写清楚：

- GOROOT :

    Golang的bin目录的位置

- GOPATH 

    Golang项目代码存放的位置

- Go Modules 
  
    Go的包管理，主要是替代GOPATH，在Go1.11之后可以在新项目使用

### 问题
我的问题主要在 每次在最外层运行main.go时会报错，错误信息为:

<code>
PS D:\code\gocache> go run .\main.go
go: warning: ignoring go.mod in $GOPATH D:\code\gocache
main.go:5:2: package gocache/day3 is not in std (D:\工作\go1.24.0\go\src\gocache\day3)
</code>

经过ds老师以及QwQ老师均未找出比较好的解决策略，最后再看过一些博客之后，将 day3里面的go.mod文件删除，
并且将cache.go 里面引用lru目录下的， 路径修改为 "gocache/day3/lru" 即可。