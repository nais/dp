package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/go-chi/chi"
	"google.golang.org/api/iterator"
	"gopkg.in/go-playground/validator.v8"

	log "github.com/sirupsen/logrus"
)

type api struct {
	client   *firestore.Client
	validate *validator.Validate
}

type AccessEntry struct {
	Subject string    `firestore:"subject" json:"subject,omitempty" validate:"required"`
	Start   time.Time `firestore:"start" json:"start,omitempty" validate:"required"`
	End     time.Time `firestore:"end" json:"end,omitempty" validate:"required"`
}

type Resource struct {
	ProjectID string `firestore:"project_id" json:"project_id,omitempty" validate:"required"`
	DatasetID string `firestore:"dataset_id" json:"dataset_id,omitempty" validate:"required"`
	Type      string `firestore:"type" json:"type,omitempty" validate:"required"`
}

type DataProduct struct {
	Name        string         `firestore:"name" json:"name,omitempty" validate:"required"`
	Description string         `firestore:"description" json:"description,omitempty" validate:"required"`
	Resource    Resource       `firestore:"resource" json:"resource,omitempty" validate:"required"`
	URI         string         `firestore:"uri" json:"uri,omitempty" validate:"required"`
	Owner       string         `firestore:"owner" json:"owner,omitempty" validate:"required"`
	Access      []*AccessEntry `firestore:"access" json:"access,omitempty" validate:"required,dive"`
}

type DataProductResponse struct {
	ID          string      `json:"id"`
	DataProduct DataProduct `json:"data_product"`
	Updated     time.Time   `json:"updated"`
	Created     time.Time   `json:"created"`
}

func (a *api) dataproducts(w http.ResponseWriter, r *http.Request) {
	dpc := a.client.Collection("dp")

	var dataproducts []DataProductResponse

	documentIterator := dpc.Documents(r.Context())
	for {
		document, err := documentIterator.Next()
		fmt.Println(document, err)
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Errorf("Query error getting dataproducts: %v", err)
		}

		var dpr DataProductResponse
		var dp DataProduct

		err = document.DataTo(&dp)

		if err != nil {
			log.Errorf("Could not deserialize document into DataProduct: %v", err)
		}

		dpr.ID = document.Ref.ID
		dpr.Updated = document.UpdateTime
		dpr.Created = document.CreateTime

		dpr.DataProduct = dp
		dataproducts = append(dataproducts, dpr)
	}

	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(dataproducts)

	if err != nil {
		log.Errorf("encoding dataproducts response: %v", err)
		respondf(w, http.StatusInternalServerError, "unable to get device config\n")
		return
	}
}

func (a *api) createDataproduct(w http.ResponseWriter, r *http.Request) {
	dpc := a.client.Collection("dp")
	var dp DataProduct

	if err := json.NewDecoder(r.Body).Decode(&dp); err != nil {
		log.Errorf("Deserializing document: %v", err)
		respondf(w, http.StatusBadRequest, "unable to deserialize document\n")
		return
	}

	if errs := a.validate.Struct(dp); errs != nil {
		log.Errorf("Validation fails: %v", errs)
		respondf(w, http.StatusBadRequest, "Validation failed: %v", errs)
	}

	documentRef, _, err := dpc.Add(r.Context(), dp)
	if err != nil {
		log.Errorf("Adding dataproduct to collection: %v", err)
		respondf(w, http.StatusInternalServerError, "unable to add dataproduct to collection\n")
		return
	}

	respondf(w, http.StatusCreated, documentRef.ID)
}

func (a *api) updateDataproduct(w http.ResponseWriter, r *http.Request) {
	dpc := a.client.Collection("dp")
	articleID := chi.URLParam(r, "productID")
	documentRef := dpc.Doc(articleID)

	var dp DataProduct

	if err := json.NewDecoder(r.Body).Decode(&dp); err != nil {
		log.Errorf("Deserializing document: %v", err)
		respondf(w, http.StatusBadRequest, "unable to deserialize document\n")
		return
	}

	var updates []firestore.Update

	if len(dp.Name) > 0 {
		updates = append(updates, firestore.Update{
			Path:  "name",
			Value: dp.Name,
		})
	}
	if len(dp.URI) > 0 {
		updates = append(updates, firestore.Update{
			Path:  "uri",
			Value: dp.URI,
		})
	}
	if len(dp.Description) > 0 {
		updates = append(updates, firestore.Update{
			Path:  "description",
			Value: dp.Description,
		})
	}
	if len(dp.Resource.Type) > 0 {
		updates = append(updates, firestore.Update{
			Path:  "resource.type",
			Value: dp.Resource.Type,
		})
	}
	if len(dp.Resource.DatasetID) > 0 {
		updates = append(updates, firestore.Update{
			Path:  "resource.dataset_id",
			Value: dp.Resource.DatasetID,
		})
	}
	if len(dp.Resource.ProjectID) > 0 {
		updates = append(updates, firestore.Update{
			Path:  "resource.project_id",
			Value: dp.Resource.ProjectID,
		})
	}
	if len(dp.Owner) > 0 {
		updates = append(updates, firestore.Update{
			Path:  "owner",
			Value: dp.Owner,
		})
	}
	_, err := documentRef.Update(r.Context(), updates)
	if err != nil {
		log.Errorf("Updating document: %v", err)
		respondf(w, http.StatusBadRequest, "unable to update document\n")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (a *api) getDataproduct(w http.ResponseWriter, r *http.Request) {
	dpc := a.client.Collection("dp")
	articleID := chi.URLParam(r, "productID")
	documentRef := dpc.Doc(articleID)

	document, err := documentRef.Get(r.Context())
	if err != nil {
		log.Errorf("Getting document: %v", err)
		respondf(w, http.StatusBadRequest, "unable to get document\n")
		return
	}

	var dpr DataProductResponse
	var dp DataProduct
	err = document.DataTo(&dp)

	if err != nil {
		log.Errorf("Could not deserialize document into DataProduct: %v", err)
	}

	dpr.ID = document.Ref.ID
	dpr.Updated = document.UpdateTime
	dpr.Created = document.CreateTime
	dpr.DataProduct = dp

	if err := json.NewEncoder(w).Encode(dpr); err != nil {
		log.Errorf("Serializing document: %v", err)
		respondf(w, http.StatusInternalServerError, "unable to serialize document\n")
		return
	}
}

func (a *api) deleteDataproduct(w http.ResponseWriter, r *http.Request) {
	dpc := a.client.Collection("dp")
	articleID := chi.URLParam(r, "productID")
	documentRef := dpc.Doc(articleID)

	if _, err := documentRef.Delete(r.Context()); err != nil {
		log.Errorf("Deleting document: %v", err)
		respondf(w, http.StatusBadRequest, "unable to delete document\n")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func respondf(w http.ResponseWriter, statusCode int, format string, args ...interface{}) {
	w.WriteHeader(statusCode)

	if _, wErr := w.Write([]byte(fmt.Sprintf(format, args...))); wErr != nil {
		log.Errorf("unable to write response: %v", wErr)
	}
}
