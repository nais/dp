package api

import (
	"cloud.google.com/go/firestore"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/nais/dp/backend/iam"
	"google.golang.org/api/iterator"

	"github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"
)

type DataProduct struct {
	Name        string               `firestore:"name" json:"name,omitempty" validate:"required"`
	Description string               `firestore:"description" json:"description,omitempty"`
	Datastore   []map[string]string  `firestore:"datastore" json:"datastore,omitempty" validate:"max=1"`
	Team        string               `firestore:"team" json:"team,omitempty" validate:"required"`
	Access      map[string]time.Time `firestore:"access" json:"access"`
}

type DataProductResponse struct {
	ID          string      `json:"id"`
	DataProduct DataProduct `json:"data_product"`
	Updated     time.Time   `json:"updated"`
	Created     time.Time   `json:"created"`
}

func (a *api) getDataproduct(w http.ResponseWriter, r *http.Request) {
	dpc := a.client.Collection(a.config.Firestore.DataproductsCollection)
	articleID := chi.URLParam(r, "productID")
	documentRef := dpc.Doc(articleID)

	document, err := documentRef.Get(r.Context())
	if err != nil {
		log.Errorf("Getting firestore document: %v", err)
		if status.Code(err) == codes.NotFound {
			respondf(w, http.StatusNotFound, "no such document\n")
		} else {
			respondf(w, http.StatusBadRequest, "unable to get document\n")
		}
		return
	}

	dpr, err := documentToProductResponse(document)
	if err != nil {
		log.Errorf("Deserializing firestore document: %v", err)
		respondf(w, http.StatusInternalServerError, "unable to deserialize document\n")
		return
	}

	if err := json.NewEncoder(w).Encode(dpr); err != nil {
		log.Errorf("Serializing dataproduct response: %v", err)
		respondf(w, http.StatusInternalServerError, "unable to serialize dataproduct response\n")
		return
	}
}

func (a *api) dataproducts(w http.ResponseWriter, r *http.Request) {
	dpc := a.client.Collection(a.config.Firestore.DataproductsCollection)
	dataproducts := make([]DataProductResponse, 0)

	iter := dpc.Documents(r.Context())
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

		dpr, err := documentToProductResponse(document)
		if err != nil {
			log.Errorf("Deserializing firestore document: %v", err)
			respondf(w, http.StatusInternalServerError, "unable to deserialize document\n")
			return
		}

		dataproducts = append(dataproducts, dpr)
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(dataproducts); err != nil {
		log.Errorf("Serializing dataproducts response: %v", err)
		respondf(w, http.StatusInternalServerError, "unable to serialize dataproduct response\n")
		return
	}
}

func (a *api) createDataproduct(w http.ResponseWriter, r *http.Request) {
	dpc := a.client.Collection(a.config.Firestore.DataproductsCollection)
	var dp DataProduct

	if err := json.NewDecoder(r.Body).Decode(&dp); err != nil {
		log.Errorf("Deserializing request document: %v", err)
		respondf(w, http.StatusBadRequest, "unable to deserialize request document\n")
		return
	}

	if errs := a.validate.Struct(dp); errs != nil {
		log.Errorf("Validation fails: %v", errs)
		respondf(w, http.StatusBadRequest, "Validation failed: %v", errs)
		return
	}

	if len(dp.Datastore) > 0 {
		if errs := ValidateDatastore(dp.Datastore[0]); errs != nil {
			log.Errorf("Validation fails: %v", errs)
			respondf(w, http.StatusBadRequest, "Validation failed: %v", errs)
			return
		}
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
	dpc := a.client.Collection(a.config.Firestore.DataproductsCollection)
	articleID := chi.URLParam(r, "productID")
	documentRef := dpc.Doc(articleID)
	document, err := documentRef.Get(r.Context())
	if err != nil {
		log.Errorf("Getting firestore document: %v", err)
		if status.Code(err) == codes.NotFound {
			respondf(w, http.StatusNotFound, "no such document\n")
		} else {
			respondf(w, http.StatusBadRequest, "unable to get firestore document\n")
		}
		return
	}

	var firebaseDp DataProduct
	if err := document.DataTo(&firebaseDp); err != nil {
		log.Errorf("Deserializing firestore document: %v", err)
		respondf(w, http.StatusInternalServerError, "unable to deserialize firestore document\n")
		return
	}

	var dp DataProduct
	if err := json.NewDecoder(r.Body).Decode(&dp); err != nil {
		log.Errorf("Deserializing request document: %v", err)
		respondf(w, http.StatusBadRequest, "unable to deserialize request document\n")
		return
	}

	updates, err := a.createUpdates(dp)
	if err != nil {
		log.Errorf("Validation fails: %v", err)
		respondf(w, http.StatusBadRequest, "Validation failed: %v", err)
		return
	}

	// In case of partial update, where an access update is made
	// without passing any datastore objects.
	if len(dp.Datastore) > 0 {
		iam.UpdateDatastoreAccess(r.Context(), dp.Datastore[0], dp.Access)
	} else {
		iam.UpdateDatastoreAccess(r.Context(), firebaseDp.Datastore[0], dp.Access)
	}

	_, err = documentRef.Update(r.Context(), updates)
	if err != nil {
		log.Errorf("Updating firestore document: %v", err)
		respondf(w, http.StatusBadRequest, "unable to update firestore document\n")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (a *api) deleteDataproduct(w http.ResponseWriter, r *http.Request) {
	dpc := a.client.Collection(a.config.Firestore.DataproductsCollection)
	articleID := chi.URLParam(r, "productID")
	documentRef := dpc.Doc(articleID)

	if _, err := documentRef.Delete(r.Context()); err != nil {
		log.Errorf("Deleting firestore document: %v", err)
		respondf(w, http.StatusBadRequest, "unable to delete firestore document\n")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (a *api) createUpdates(dp DataProduct) ([]firestore.Update, error) {
	var updates []firestore.Update

	if len(dp.Name) > 0 {
		updates = append(updates, firestore.Update{
			Path:  "name",
			Value: dp.Name,
		})
	}
	if len(dp.Description) > 0 {
		updates = append(updates, firestore.Update{
			Path:  "description",
			Value: dp.Description,
		})
	}
	if len(dp.Datastore) > 0 {
		if errs := ValidateDatastore(dp.Datastore[0]); errs != nil {
			return nil, errs
		}
		updates = append(updates, firestore.Update{
			Path:  "datastore",
			Value: dp.Datastore,
		})
	}
	if len(dp.Team) > 0 {
		updates = append(updates, firestore.Update{
			Path:  "owner",
			Value: dp.Team,
		})
	}

	return updates, nil
}

func ValidateDatastore(store map[string]string) error {
	datastoreType := store["type"]
	if len(datastoreType) == 0 {
		return fmt.Errorf("no type defined")
	}

	switch datastoreType {
	case iam.BucketType:
		return hasKeys(store, "project_id", "bucket_id")
	case iam.BigQueryType:
		return hasKeys(store, "dataset_id", "project_id", "resource_id")
	}

	return fmt.Errorf("unknown datastore type: %v", datastoreType)
}

func hasKeys(m map[string]string, keys ...string) error {
	for _, k := range keys {
		if _, found := m[k]; !found {
			return fmt.Errorf("missing key: %v", k)
		}
	}
	return nil
}

func documentToProductResponse(d *firestore.DocumentSnapshot) (DataProductResponse, error) {
	var dpr DataProductResponse
	var dp DataProduct

	if err := d.DataTo(&dp); err != nil {
		return dpr, err
	}

	if dp.Access == nil {
		dp.Access = make(map[string]time.Time)
	}
	dpr.ID = d.Ref.ID
	dpr.Updated = d.UpdateTime
	dpr.Created = d.CreateTime
	dpr.DataProduct = dp

	return dpr, nil
}
