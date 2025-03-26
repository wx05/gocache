package final

import (
	"gopkg.in/yaml.v3"
	"os"
)

/*
配置文件读取
*/

type PretaskConfig struct {
	DataType string `yaml:"data_type"`
	FilePath string `yaml:"file_path"`
	Dsn      string `yaml:"dsn"`
}

type Config struct {
	Server  ServerConfig  `yaml:"server"`
	PreTask PretaskConfig `yaml:"preTask"`
}

type ServerConfig struct {
	MaxCacheBytes int64  `yaml:"max_cache_bytes"`
	BasePath      string `yaml:"base_path"`
	Replicas      int    `yaml:"replicas"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
