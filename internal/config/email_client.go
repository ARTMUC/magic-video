package config

type BrevoEmailClientConfig struct {
	apiKey      string
	senderName  string
	senderEmail string
}

func (b BrevoEmailClientConfig) SenderName() string {
	return b.senderName
}

func (b BrevoEmailClientConfig) SenderEmail() string {
	return b.senderEmail
}

func (b BrevoEmailClientConfig) ApiKey() string {
	return b.apiKey
}
