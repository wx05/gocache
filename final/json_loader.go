package final

import (
	"encoding/json"
	"os"
)

type JsonLoader struct{}

type GroupData struct {
	Group string     `json:"group"`
	Data  []KeyValue `json:"data"`
}

type KeyValue struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

func (j *JsonLoader) Load(c *Config) ([]map[string]interface{}, error) {
	data, err := os.ReadFile(c.PreTask.FilePath)
	if err != nil {
		return nil, err
	}

	var jsonData []GroupData
	err = json.Unmarshal(data, &jsonData)
	if err != nil {
		return nil, err
	}

	return make([]map[string]interface{}, 0), nil
}
