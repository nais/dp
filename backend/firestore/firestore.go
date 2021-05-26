package firestore

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"time"
)

type Firestore struct {
	dataproducts *firestore.CollectionRef
}

type Dataproduct struct {
	Name        string               `firestore:"name"`
	Description string               `firestore:"description"`
	Datastore   []map[string]string  `firestore:"datastore"`
	Team        string               `firestore:"team"`
	Access      map[string]time.Time `firestore:"access"`
}

func New(ctx context.Context, googleProjectID, dataproductCollection, accessCollection string) (*Firestore, error) {
	client, err := firestore.NewClient(ctx, googleProjectID)
	if err != nil {
		return nil, fmt.Errorf("initializing firestore client: %v", err)
	}
	return &Firestore{
		dataproducts: client.Collection(dataproductCollection),
	}, nil
}

func (f *Firestore) CreateDataproduct(ctx context.Context, dp Dataproduct) (string, error) {
	ref, _, err := f.dataproducts.Add(ctx, dp)
	if err != nil {
		log.Errorf("Adding dataproduct to collection: %v", err)
		return "", fmt.Errorf("adding dataproduct to collection: %w", err)
	}
	return ref.ID, nil
}
