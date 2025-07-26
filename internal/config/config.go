package config

type Config struct {
	serverConfig           ServerConfig
	encryptionConfig       EncryptionConfig
	brevoEmailClientConfig BrevoEmailClientConfig
	przelewy24ClientConfig Przelewy24ClientConfig
	sessionConfig          SessionConfig
}

func (c Config) Przelewy24ClientConfig() Przelewy24ClientConfig {
	return c.przelewy24ClientConfig
}

func (c Config) SessionConfig() SessionConfig {
	return c.sessionConfig
}

func (c Config) EncryptionConfig() EncryptionConfig {
	return c.encryptionConfig
}

func (c Config) BrevoEmailClientConfig() BrevoEmailClientConfig {
	return c.brevoEmailClientConfig
}

func (c Config) ServerConfig() ServerConfig {
	return c.serverConfig
}
