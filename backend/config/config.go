package config

type Config struct {
	BindAddress string
	LogLevel    string
	DevMode     bool
	Firestore   FirestoreConfig
	OAuth2      OAuth2Config
	TeamsURL    string
	TeamsToken  string
	Hostname    string
	State       string
}

type FirestoreConfig struct {
	GoogleProjectID         string
	DataproductsCollection  string
	AccessUpdatesCollection string
}

type OAuth2Config struct {
	ClientID     string
	ClientSecret string
	TenantID     string
}

func DefaultConfig() Config {
	return Config{
		BindAddress: ":8080",
		LogLevel:    "info",
		Firestore: FirestoreConfig{
			GoogleProjectID: "aura-dev-d9f5",
		},
		TeamsURL: "https://raw.githubusercontent.com/navikt/teams/main/teams.json",
	}
}
