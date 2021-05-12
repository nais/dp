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

type AccessEntry struct {
	Subject string    `firestore:"subject" json:"subject,omitempty" validate:"required"`
	Start   time.Time `firestore:"start" json:"start,omitempty" validate:"required"`
	End     time.Time `firestore:"end" json:"end,omitempty" validate:"required"`
}

type DataProduct struct {
	Name        string              `firestore:"name" json:"name,omitempty" validate:"required"`
	Description string              `firestore:"description" json:"description,omitempty"`
	Datastore   []map[string]string `firestore:"datastore" json:"datastore,omitempty" validate:"max=1"`
	Owner       string              `firestore:"owner" json:"owner,omitempty" validate:"required"`
	Access      []*AccessEntry      `firestore:"access" json:"access" validate:"dive"`
}

type DataProductResponse struct {
	ID          string      `json:"id"`
	DataProduct DataProduct `json:"data_product"`
	Updated     time.Time   `json:"updated"`
	Created     time.Time   `json:"created"`
}
