package firestore

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/iterator"
	"time"
)

type Firestore struct {
	dataproducts *firestore.CollectionRef
}

type Dataproduct struct {
	Name        string               `firestore:"name" json:"name,omitempty"`
	Description string               `firestore:"description" json:"description,omitempty"`
	Datastore   []map[string]string  `firestore:"datastore" json:"datastore,omitempty"`
	Team        string               `firestore:"team" json:"team,omitempty"`
	Access      map[string]time.Time `firestore:"access" json:"access,omitempty"`
}

type DataproductResponse struct {
	ID          string      `json:"id"`
	Dataproduct Dataproduct `json:"data_product"`
	Updated     time.Time   `json:"updated"`
	Created     time.Time   `json:"created"`
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

func (f *Firestore) GetDataproduct(ctx context.Context, id string) (*DataproductResponse, error) {
	doc, err := f.dataproducts.Doc(id).Get(ctx)
	if err != nil {
		log.Errorf("Getting dataproduct from collection: %v", err)
		return nil, fmt.Errorf("getting dataproduct from collection: %w", err)
	}

	return toResponse(doc)
}

func (f *Firestore) GetDataproducts(ctx context.Context) ([]*DataproductResponse, error) {
	var dataproducts []*DataproductResponse

	iter := f.dataproducts.Documents(ctx)
	defer iter.Stop()
	for {
		document, err := iter.Next()

		if err == iterator.Done {
			break
		}

		if err != nil {
			log.Errorf("Iterating documents: %v", err)
			break
		}

		dpr, err := toResponse(document)
		if err != nil {
			log.Errorf("Creating DataproductResponse: %v", err)
			return nil, fmt.Errorf("creating DataproductResponse: %w", err)
		}

		dataproducts = append(dataproducts, dpr)
	}

	return dataproducts, nil
}

func toResponse(document *firestore.DocumentSnapshot) (*DataproductResponse, error) {
	var dp Dataproduct

	if err := document.DataTo(&dp); err != nil {
		return nil, fmt.Errorf("populating fields in dataproduct struct: %w", err)
	}

	var dpr DataproductResponse
	dpr.Dataproduct = dp
	dpr.ID = document.Ref.ID
	dpr.Updated = document.UpdateTime
	dpr.Created = document.CreateTime

	return &dpr, nil
}
