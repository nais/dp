package config

type Config struct {
	BindAddress              string
	LogLevel                 string
	DevMode                  bool
	FirestoreGoogleProjectId string
	FirestoreCollection      string
	OAuth2ClientID           string
	OAuth2ClientSecret       string
	OAuth2TenantID           string
}

func DefaultConfig() Config {
	return Config{
		BindAddress:              ":8080",
		LogLevel:                 "info",
		FirestoreGoogleProjectId: "aura-dev-d9f5",
	}
}
