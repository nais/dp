package api

import (
	"cloud.google.com/go/firestore"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"google.golang.org/api/iterator"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

type api struct {
	client *firestore.Client
}

type DataProductRequest struct {
	Name        string `firestore:"name" json:"name,omitempty"`
	Description string `firestore:"description" json:"description,omitempty"`
	Type        string `firestore:"type" json:"type,omitempty"`
	URI         string `firestore:"uri" json:"uri,omitempty"`
}

type DataProductResponse struct {
	ID    		string 	  `json:"id"`
	Name        string    `firestore:"name" json:"name,omitempty"`
	Description string    `firestore:"description" json:"description,omitempty"`
	Type        string    `firestore:"type" json:"type,omitempty"`
	URI         string 	  `firestore:"uri" json:"uri,omitempty"`
	Updated     time.Time `json:"updated"`
	Created     time.Time `json:"created"`
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

		var dp DataProductResponse
		err = document.DataTo(&dp)

		if err != nil {
			log.Errorf("Could not deserialize document into DataProductRequest: %v", err)
		}

		dp.ID = document.Ref.ID
		dp.Updated = document.UpdateTime
		dp.Created = document.CreateTime

		dataproducts = append(dataproducts, dp)
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
	var dp DataProductRequest

	if err := json.NewDecoder(r.Body).Decode(&dp); err != nil {
		log.Errorf("Deserializing document: %v", err)
		respondf(w, http.StatusBadRequest, "unable to deserialize document\n")
		return
	}

	if len(dp.Name) == 0 {
		log.Errorf("Missing required field: name")
		respondf(w, http.StatusBadRequest, "Missing required field: name")
		return
	}
	if len(dp.URI) == 0 {
		log.Errorf("Missing required field: uri")
		respondf(w, http.StatusBadRequest, "Missing required field: uri")
		return
	}
	if len(dp.Description) == 0 {
		log.Errorf("Missing required field: description")
		respondf(w, http.StatusBadRequest, "Missing required field: description")
		return
	}
	if len(dp.Type) == 0 {
		log.Errorf("Missing required field: type")
		respondf(w, http.StatusBadRequest, "Missing required field: type")
		return
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

	var dp DataProductRequest

	if err := json.NewDecoder(r.Body).Decode(&dp); err != nil {
		log.Errorf("Deserializing document: %v", err)
		respondf(w, http.StatusBadRequest, "unable to deserialize document\n")
		return
	}

	var updates []firestore.Update

	if dp.Name != "" {
		updates = append(updates, firestore.Update{
			Path:      "name",
			Value:     dp.Name,
		})
	}
	if dp.URI != "" {
		updates = append(updates, firestore.Update{
			Path:      "uri",
			Value:     dp.URI,
		})
	}
	if dp.Description != "" {
		updates = append(updates, firestore.Update{
			Path:      "description",
			Value:     dp.Description,
		})
	}
	if dp.Type != "" {
		updates = append(updates, firestore.Update{
			Path:      "type",
			Value:     dp.Type,
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

	var dp DataProductResponse

	err = document.DataTo(&dp)

	if err != nil {
		log.Errorf("Could not deserialize document into DataProductRequest: %v", err)
	}

	dp.ID = document.Ref.ID
	dp.Updated = document.UpdateTime
	dp.Created = document.CreateTime

	if err := json.NewEncoder(w).Encode(dp); err != nil {
		log.Errorf("Serializing document: %v", err)
		respondf(w, http.StatusInternalServerError, "unable to serialize document\n")
		return
	}
}

func respondf(w http.ResponseWriter, statusCode int, format string, args ...interface{}) {
	w.WriteHeader(statusCode)

	if _, wErr := w.Write([]byte(fmt.Sprintf(format, args...))); wErr != nil {
		log.Errorf("unable to write response: %v", wErr)
	}
}
