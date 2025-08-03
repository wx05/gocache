package final

import "fmt"

type LoadFactory struct{}

func (l *LoadFactory) CreateLoader(c *Config) (DataLoader, error) {
	switch c.PreTask.DataType {
	case "json":
		return &JsonLoader{}, nil
	case "file":
		return &FileLoader{}, nil
	case "db":
		return &DbLoader{}, nil
	default:
		return nil, fmt.Errorf("unsupport loader, %s", c.PreTask.DataType)
	}
}
