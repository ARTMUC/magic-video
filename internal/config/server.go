package config

type ServerConfig struct {
	port       string
	apiUrl     string
	websiteUrl string
}

func (c ServerConfig) WebsiteUrl() string {
	return c.websiteUrl
}

func (c ServerConfig) ApiUrl() string {
	return c.apiUrl
}

func (c ServerConfig) Port() string {
	return c.port
}
