package final

type DbLoader struct{}

func (j *DbLoader) Load(c *Config, peer *HTTPPool) ([]*Group, error) {
	return []*Group{}, nil
}
