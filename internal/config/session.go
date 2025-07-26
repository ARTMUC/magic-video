package config

type ServerConfig struct {
	port   string
	apiUrl string
}

func (c ServerConfig) ApiUrl() string {
	return c.apiUrl
}

func (c ServerConfig) Port() string {
	return c.port
}
