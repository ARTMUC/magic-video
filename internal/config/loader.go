package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func LoadConfig(isLocal bool) (*Config, error) {
	if isLocal {
		err := godotenv.Load(".env")
		if err != nil {
			return nil, fmt.Errorf("can't load env variables: %w", err)
		}
	}

	config := &Config{
		serverConfig: ServerConfig{
			port:       os.Getenv("PORT"),
			apiUrl:     os.Getenv("API_URL"),
			websiteUrl: os.Getenv("WEBSITE_URL"),
		},
		encryptionConfig: EncryptionConfig{
			encryptionKey: os.Getenv("ENCRYPTION_KEY"),
		},
		brevoEmailClientConfig: BrevoEmailClientConfig{
			apiKey:      os.Getenv("BREVO_API_KEY"),
			senderName:  os.Getenv("BREVO_SENDER_NAME"),
			senderEmail: os.Getenv("BREVO_SENDER_EMAIL"),
		},
		przelewy24ClientConfig: Przelewy24ClientConfig{
			merchatID:   os.Getenv("P24_MERCHANT_ID"),
			posID:       os.Getenv("P24_POS_ID"),
			apiKey:      os.Getenv("P24_API_KEY"),
			crc:         os.Getenv("P24_CRC"),
			salt:        os.Getenv("P24_SALT"),
			environment: os.Getenv("P24_ENVIRONMENT"),
			returnUrl:   os.Getenv("P24_RETURN_URL"),
			webhookUrl:  os.Getenv("P24_WEBHOOK_URL"),
		},
		sessionConfig: SessionConfig{
			jwtTokenSecret:        os.Getenv("JWT_TOKEN_SECRET"),
			jwtRefreshTokenSecret: os.Getenv("JWT_REFRESH_TOKEN_SECRET"),
			jwtTokenExpiry:        toTimeDuration(os.Getenv("JWT_TOKEN_EXPIRY"), 15*60),
			jwtRefreshTokenExpiry: toTimeDuration(os.Getenv("JWT_REFRESH_TOKEN_EXPIRY"), 24*60*60),
		},
	}

	return config, nil
}
