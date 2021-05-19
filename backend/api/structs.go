package api

import (
	"time"

	"cloud.google.com/go/firestore"
	"github.com/nais/dp/backend/config"
	"gopkg.in/go-playground/validator.v9"
)

type api struct {
	client    *firestore.Client
	validate  *validator.Validate
	config    config.Config
	teamUUIDs map[string]string
}

type DataProduct struct {
	Name        string               `firestore:"name" json:"name,omitempty" validate:"required"`
	Description string               `firestore:"description" json:"description,omitempty"`
	Datastore   []map[string]string  `firestore:"datastore" json:"datastore,omitempty" validate:"max=1"`
	Owner       string               `firestore:"owner" json:"owner,omitempty" validate:"required"`
	Access      map[string]time.Time `firestore:"access" json:"access"`
}

type DataProductResponse struct {
	ID          string      `json:"id"`
	DataProduct DataProduct `json:"data_product"`
	Updated     time.Time   `json:"updated"`
	Created     time.Time   `json:"created"`
}
