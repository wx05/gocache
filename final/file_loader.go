package final

type FileLoader struct{}

func (j *FileLoader) Load(c *Config) ([]map[string]interface{}, error) {
	return make([]map[string]interface{}, 0), nil
}
