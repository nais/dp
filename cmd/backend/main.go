package main

import (
	"context"
	"fmt"
	"github.com/nais/dp/backend/config"
	"github.com/nais/dp/backend/middleware"
	"log"
	"net/http"
	"os"

	firestore "cloud.google.com/go/firestore"
	"github.com/nais/dp/backend/api"
	flag "github.com/spf13/pflag"
)

var cfg = config.DefaultConfig()

func init() {
	flag.StringVar(&cfg.BindAddress, "bind-address", cfg.BindAddress, "Bind address")
	flag.StringVar(&cfg.OAuth2.ClientID, "oauth2-client-id", os.Getenv("OAUTH2_CLIENT_ID"), "OAuth2 client ID")
	flag.StringVar(&cfg.OAuth2.ClientSecret, "oauth2-client-secret", os.Getenv("OAUTH2_CLIENT_SECRET"), "OAuth2 client secret")
	flag.BoolVar(&cfg.DevMode, "development-mode", cfg.DevMode, "Run in development mode")
	flag.Parse()
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client, err := firestore.NewClient(ctx, "aura-dev-d9f5")
	if err != nil {
		log.Fatalf("Initializing firestore client: %v", err)
	}

	router := api.New(client, jwtValidatorMiddleware(cfg))
	fmt.Println("running @", "localhost:8080")
	fmt.Println(http.ListenAndServe("localhost:8080", router))
}

func jwtValidatorMiddleware(c config.Config) func(http.Handler) http.Handler {
	if c.DevMode {
		return middleware.MockJWTValidatorMiddleware()
	}
	return middleware.JWTValidatorMiddleware(cfg.OAuth2)
}