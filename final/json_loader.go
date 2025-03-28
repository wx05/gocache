package final

import (
	"encoding/json"
	"os"
)

type JsonLoader struct{}

type GroupData struct {
	Group string            `json:"group"`
	Data  map[string]string `json:"data"`
}

func (j *JsonLoader) Load(c *Config, peer *HTTPPool) ([]*Group, error) {
	data, err := os.ReadFile(c.PreTask.FilePath)
	if err != nil {
		return nil, err
	}

	var jsonData []GroupData
	err = json.Unmarshal(data, &jsonData)
	if err != nil {
		return nil, err
	}

	var groupList []*Group
	for _, v := range jsonData {
		group := CreateGroup(c, v)
		groupList = append(groupList, group)

		//仅仅只在数据被选择的节点进行初始化
		for k1, v1 := range v.Data {
			if ok := peer.IsSelf(k1); ok == true && group.name == v.Group {
				group.mainCache.add(k1, ByteView{b: []byte(v1)})
			}
		}
	}

	return groupList, nil
}
