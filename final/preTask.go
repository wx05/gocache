package final

import "log"

/*
主要为数据预加载的实现
*/

type DataLoader interface {
	Load(c *Config) ([]map[string]interface{}, error)
}

type PreTask struct {
}

func (p *PreTask) Run(config *Config) {

	factory := LoadFactory{}
	loader, err := factory.CreateLoader(config)
	if err != nil {
		log.Panicf("create loader error, %s", err.Error())
		return
	}

	_, err = loader.Load(config)
	if err != nil {
		log.Panicf("preload data error, %s", err.Error())
		return
	}

	//写入缓存,这里的Group需要和 后面的 链接在一起
}
