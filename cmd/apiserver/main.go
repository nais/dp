package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	firestore "cloud.google.com/go/firestore"
	"github.com/nais/dp/apiserver/api"
	flag "github.com/spf13/pflag"
	"gopkg.in/go-playground/validator.v8"
)

var validate *validator.Validate

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
	config := &validator.Config{TagName: "validate"}
	validate = validator.New(config)
	router := api.New(client, validate)
	fmt.Println("running @", "localhost:8080")
	fmt.Println(http.ListenAndServe("localhost:8080", router))
}
