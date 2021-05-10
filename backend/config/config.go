package config

import (
	"github.com/nais/dp/backend/auth"
)

type Config struct {
	BindAddress              string
	OAuth2                   auth.Google
	LogLevel                 string
	DevMode                  bool
	FirestoreGoogleProjectId string
	FirestoreCollection      string
}

func DefaultConfig() Config {
	return Config{
		OAuth2: auth.Google{
			DiscoveryURL: "https://accounts.google.com/.well-known/openid-configuration",
		},
		BindAddress:              ":8080",
		LogLevel:                 "info",
		FirestoreGoogleProjectId: "aura-dev-d9f5",
	}
}
