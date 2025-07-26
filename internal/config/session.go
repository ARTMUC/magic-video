package config

import "time"

type SessionConfig struct {
	jwtTokenSecret        string
	jwtRefreshTokenSecret string
	jwtTokenExpiry        time.Duration
	jwtRefreshTokenExpiry time.Duration
}

func (s SessionConfig) JwtTokenSecret() string {
	return s.jwtTokenSecret
}

func (s SessionConfig) JwtRefreshTokenSecret() string {
	return s.jwtRefreshTokenSecret
}

func (s SessionConfig) JwtTokenExpiry() time.Duration {
	return s.jwtTokenExpiry
}

func (s SessionConfig) JwtRefreshTokenExpiry() time.Duration {
	return s.jwtRefreshTokenExpiry
}
