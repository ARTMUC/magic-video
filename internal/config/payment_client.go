package config

type Przelewy24ClientConfig struct {
	merchatID   string
	posID       string
	apiKey      string
	crc         string
	salt        string
	environment string
	returnUrl   string
	webhookUrl  string
}

func (p Przelewy24ClientConfig) WebhookUrl() string {
	return p.webhookUrl
}

func (p Przelewy24ClientConfig) ReturnUrl() string {
	return p.returnUrl
}

func (p Przelewy24ClientConfig) MerchatID() string {
	return p.merchatID
}

func (p Przelewy24ClientConfig) PosID() string {
	return p.posID
}

func (p Przelewy24ClientConfig) ApiKey() string {
	return p.apiKey
}

func (p Przelewy24ClientConfig) Crc() string {
	return p.crc
}

func (p Przelewy24ClientConfig) Environment() string {
	return p.environment
}

func (p Przelewy24ClientConfig) Salt() string {
	return p.salt
}
