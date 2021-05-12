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
	TeamsURL                 string
	TeamsToken               string
}

func DefaultConfig() Config {
	return Config{
		BindAddress:              ":8080",
		LogLevel:                 "info",
		FirestoreGoogleProjectId: "aura-dev-d9f5",
		TeamsURL:                 "https://raw.githubusercontent.com/navikt/teams/main/teams.json",
	}
}
