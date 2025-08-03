package final

type FileLoader struct{}

func (j *FileLoader) Load(c *Config, peer *HTTPPool) ([]*Group, error) {
	return []*Group{}, nil
}
