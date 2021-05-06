package main

import (
	"context"
	"fmt"
	"github.com/nais/dp/apiserver/auth"
	"log"
	"net/http"

	firestore "cloud.google.com/go/firestore"
	"github.com/nais/dp/apiserver/api"
	flag "github.com/spf13/pflag"
	"gopkg.in/go-playground/validator.v9"
)

var inputValidator *validator.Validate

func init() {
	//flag.StringVar(&cfg.BindAddress, "bind-address", cfg.BindAddress, "Bind address")
	flag.Parse()
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client, err := firestore.NewClient(ctx, "aura-dev-d9f5")
	if err != nil {
		log.Fatalf("Initializing firestore client: %v", err)
	}

	jwtValidator, err := auth.CreateJWTValidator(auth.Google{
		DiscoveryURL: "https://accounts.google.com/.well-known/openid-configuration",
		ClientID:     "854073996265-riks3c6p36oh3ijgef8tvlk3367ab9sq.apps.googleusercontent.com",
	})

	inputValidator = validator.New()
	router := api.New(client, inputValidator, jwtValidator)
	fmt.Println("running @", "localhost:8080")
	fmt.Println(http.ListenAndServe("localhost:8080", router))
}
