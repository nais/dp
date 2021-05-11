package main

import (
	"context"
	"fmt"
	"github.com/nais/dp/backend/auth"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/nais/dp/backend/config"

	firestore "cloud.google.com/go/firestore"
	"github.com/nais/dp/backend/api"
	flag "github.com/spf13/pflag"
)

var cfg = config.DefaultConfig()

const TeamsUpdateFrequency = 5 * time.Minute

func init() {
	flag.StringVar(&cfg.BindAddress, "bind-address", cfg.BindAddress, "Bind address")
	flag.StringVar(&cfg.OAuth2ClientID, "oauth2-client-id", os.Getenv("AZURE_APP_CLIENT_ID"), "OAuth2 client ID")
	flag.StringVar(&cfg.OAuth2ClientSecret, "oauth2-client-secret", os.Getenv("AZURE_APP_CLIENT_SECRET"), "OAuth2 client secret")
	flag.StringVar(&cfg.OAuth2TenantID, "oauth2-tenant-id", os.Getenv("AZURE_APP_TENANT_ID"), "Azure tenant id")
	flag.StringVar(&cfg.FirestoreGoogleProjectId, "firestore-google-project-id", os.Getenv("FIRESTORE_GOOGLE_PROJECT_ID"), "Firestore Google project ID")
	flag.StringVar(&cfg.FirestoreCollection, "firestore-collection", os.Getenv("FIRESTORE_COLLECTION"), "Firestore collection name")
	flag.StringVar(&cfg.TeamsURL, "teams-url", os.Getenv("TEAMS_URL"), "URL for json containing teams and UUIDs")
	flag.StringVar(&cfg.TeamsToken, "teams-token", os.Getenv("TEAMS_TOKEN"), "Token for accessing teams json")
	flag.BoolVar(&cfg.DevMode, "development-mode", cfg.DevMode, "Run in development mode")
	flag.Parse()
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client, err := firestore.NewClient(ctx, cfg.FirestoreGoogleProjectId)
	if err != nil {
		log.Fatalf("Initializing firestore client: %v", err)
	}

	teamUUIDs := make(map[string]string)
	go auth.UpdateTeams(ctx, teamUUIDs, cfg.TeamsURL, cfg.TeamsToken, TeamsUpdateFrequency)

	api := api.New(client, cfg, teamUUIDs)
	fmt.Println("running @", cfg.BindAddress)
	fmt.Println(http.ListenAndServe(cfg.BindAddress, api))
}
