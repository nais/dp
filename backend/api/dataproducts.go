package api

import (
	"encoding/json"
	"fmt"
	firestore2 "github.com/nais/dp/backend/firestore"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/go-chi/chi"
	"github.com/nais/dp/backend/iam"
	log "github.com/sirupsen/logrus"
)

type DataProduct struct {
	DataProductInput
	Access map[string]time.Time `firestore:"access" json:"access"`
}

type DataProductInput struct {
	Name        string              `firestore:"name" json:"name,omitempty" validate:"required"`
	Description string              `firestore:"description" json:"description,omitempty"`
	Datastore   []map[string]string `firestore:"datastore" json:"datastore,omitempty" validate:"max=1"`
	Team        string              `firestore:"team" json:"team,omitempty" validate:"required"`
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

	requesterTeams := r.Context().Value("teams").([]string)

	if !contains(requesterTeams, firebaseDp.Team) {
		log.Errorf("updateDataproduct: Unauthorized to make changes to firestore document")
		respondf(w, http.StatusUnauthorized, "unauthorized\n")
		return
	}

	var dpi DataProductInput
	if err := json.NewDecoder(r.Body).Decode(&dpi); err != nil {
		log.Errorf("Deserializing request document: %v", err)
		respondf(w, http.StatusBadRequest, "unable to deserialize request document\n")
		return
	}

	updates, err := a.createUpdates(dpi, firebaseDp.Team, firebaseDp.Access)
	if err != nil {
		log.Errorf("Validation fails: %v", err)
		respondf(w, http.StatusBadRequest, "Validation failed: %v", err)
		return
	}

	_, err = documentRef.Update(r.Context(), updates)
	if err != nil {
		log.Errorf("Updating firestore document: %v", err)
		respondf(w, http.StatusBadRequest, "unable to update firestore document\n")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (a *api) dataproducts(w http.ResponseWriter, r *http.Request) {
	dataproducts, err := a.firestore.GetDataproducts(r.Context())
	if err != nil {
		log.Errorf("Getting dataproducts: %v", err)
		respondf(w, http.StatusInternalServerError, "unable to get dataproducts")
		return
	}

	if err := json.NewEncoder(w).Encode(dataproducts); err != nil {
		log.Errorf("Serializing dataproducts response: %v", err)
		respondf(w, http.StatusInternalServerError, "unable to serialize dataproduct response\n")
		return
	}
}

func (a *api) createDataproduct(w http.ResponseWriter, r *http.Request) {
	var dpi DataProductInput
	var dp firestore2.Dataproduct

	if err := json.NewDecoder(r.Body).Decode(&dpi); err != nil {
		log.Errorf("Deserializing request document: %v", err)
		respondf(w, http.StatusBadRequest, "unable to deserialize request document\n")
		return
	}

	if errs := a.validate.Struct(dpi); errs != nil {
		log.Errorf("Validation fails: %v", errs)
		respondf(w, http.StatusBadRequest, "Validation failed: %v", errs)
		return
	}

	if len(dpi.Datastore) > 0 {
		if errs := ValidateDatastore(dp.Datastore[0]); errs != nil {
			log.Errorf("Validation fails: %v", errs)
			respondf(w, http.StatusBadRequest, "Validation failed: %v", errs)
			return
		}
	}

	dp.Access = make(map[string]time.Time)
	dp.Access[fmt.Sprintf("group:%v@nav.no", dpi.Team)] = time.Time{} // gives infinite access to the owners (team) of the dataproduct
	dp.Datastore = dpi.Datastore
	dp.Team = dpi.Team
	dp.Name = dpi.Name
	dp.Description = dpi.Description

	id, err := a.firestore.CreateDataproduct(r.Context(), dp)

	if err != nil {
		respondf(w, http.StatusInternalServerError, "unable to create dataproduct\n")
		return
	}

	respondf(w, http.StatusCreated, id)
}

func (a *api) getDataproduct(w http.ResponseWriter, r *http.Request) {
	dataproduct, err := a.firestore.GetDataproduct(r.Context(), chi.URLParam(r, "productID"))

	if err != nil {
		log.Errorf("Getting firestore document: %v", err)
		if status.Code(err) == codes.NotFound {
			respondf(w, http.StatusNotFound, "not found\n")
		} else {
			respondf(w, http.StatusBadRequest, "unable to get document\n")
		}
	}

	if err := json.NewEncoder(w).Encode(dataproduct); err != nil {
		log.Errorf("Serializing dataproduct response: %v", err)
		respondf(w, http.StatusInternalServerError, "unable to serialize dataproduct response\n")
		return
	}
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

func (a *api) createUpdates(dp DataProductInput, currentTeam string, access map[string]time.Time) ([]firestore.Update, error) {
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
			Path:  "team",
			Value: dp.Team,
		})

		delete(access, fmt.Sprintf("group:%v@nav.no", currentTeam))
		access[fmt.Sprintf("group:%v@nav.no", dp.Team)] = time.Time{}
		updates = append(updates, firestore.Update{
			Path:  "access",
			Value: access,
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
