package final

type DbLoader struct{}

func (j *DbLoader) Load(c *Config) ([]map[string]interface{}, error) {
	return make([]map[string]interface{}, 0), nil
}
