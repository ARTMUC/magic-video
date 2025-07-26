package config

type EncryptionConfig struct {
	encryptionKey string
}

func (e EncryptionConfig) EncryptionKey() string {
	return e.encryptionKey
}
