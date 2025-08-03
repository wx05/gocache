package final

import "log"

/*
主要为数据预加载的实现
*/

// 数据初始化存储
var Db map[string]string

type DataLoader interface {
	Load(c *Config, peer *HTTPPool) ([]*Group, error)
}

type PreTask struct {
}

func (p *PreTask) Run(config *Config, peer *HTTPPool) []*Group {

	//工厂模式load数据
	factory := LoadFactory{}
	loader, err := factory.CreateLoader(config)
	if err != nil {
		log.Panicf("create loader error, %s", err.Error())
		return nil
	}

	groups, err := loader.Load(config, peer)
	if err != nil {
		log.Panicf("preload data error, %s", err.Error())
		return nil
	}

	//如果为空,则走默认库
	if groups == nil {
		group := NewGroup(config.Server.DefalutGroup, config.Server.MaxCacheBytes, nil)
		groups = append(groups, group)
	}
	return groups
}
