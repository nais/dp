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
}

func DefaultConfig() Config {
	return Config{
		OAuth2: auth.Google{
			DiscoveryURL: "https://accounts.google.com/.well-known/openid-configuration",
			ClientID:     "854073996265-riks3c6p36oh3ijgef8tvlk3367ab9sq.apps.googleusercontent.com",
		},
		BindAddress:              ":8080",
		LogLevel:                 "info",
		FirestoreGoogleProjectId: "aura-dev-d9f5",
	}
}
