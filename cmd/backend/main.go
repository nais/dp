package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/nais/dp/backend/config"

	firestore "cloud.google.com/go/firestore"
	"github.com/nais/dp/backend/api"
	flag "github.com/spf13/pflag"
)

var cfg = config.DefaultConfig()

func init() {
	flag.StringVar(&cfg.BindAddress, "bind-address", cfg.BindAddress, "Bind address")
	flag.StringVar(&cfg.OAuth2ClientID, "oauth2-client-id", os.Getenv("AZURE_APP_CLIENT_ID"), "OAuth2 client ID")
	flag.StringVar(&cfg.OAuth2ClientSecret, "oauth2-client-secret", os.Getenv("AZURE_APP_CLIENT_SECRET"), "OAuth2 client secret")
	flag.StringVar(&cfg.OAuth2TenantID, "oauth2-tenant-id", os.Getenv("AZURE_APP_TENANT_ID"), "Azure tenant id")
	flag.StringVar(&cfg.FirestoreGoogleProjectId, "firestore-google-project-id", os.Getenv("FIRESTORE_GOOGLE_PROJECT_ID"), "Firestore Google project ID")
	flag.StringVar(&cfg.FirestoreCollection, "firestore-collection", os.Getenv("FIRESTORE_COLLECTION"), "Firestore collection name")
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

	router := api.New(client, cfg)
	fmt.Println("running @", cfg.BindAddress)
	fmt.Println(http.ListenAndServe(cfg.BindAddress, router))
}
