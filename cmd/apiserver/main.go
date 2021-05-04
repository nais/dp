package main

import (
	firestore "cloud.google.com/go/firestore"
	"context"
	"fmt"
	"github.com/nais/dp/apiserver/api"
	flag "github.com/spf13/pflag"
	"log"
)

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

	dpc := client.Collection("dp")

	fmt.Println("dpc: ", dpc)

	sampleDP := api.DataProduct{Name: "container resource usage", Description: "beskrivelse", URI: "https://uri.com"}
	doc, wr, err := dpc.Add(ctx, sampleDP)
	if err != nil {
		log.Printf("Adding sample doc: %v\n", err)
	}

	fmt.Println(doc, wr)

	//router := api.New()
	//
	//fmt.Println("running @", "localhost:8080")
	//fmt.Println(http.ListenAndServe("localhost:8080", router))
}
