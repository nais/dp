package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"cloud.google.com/go/firestore"
	"github.com/go-chi/chi"
	"github.com/nais/dp/backend/iam"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/iterator"
)

const (
	Delete = "delete"
	Grant  = "grant"
)

type AccessResponse struct {
	ID      string               `json:"id"`
	Access  map[string]time.Time `json:"access"`
	Updated time.Time            `json:"updated"`
	Created time.Time            `json:"created"`
}

type AccessSubject struct {
	Subject string    `json:"subject" validate:"required"`
	Expires time.Time `json:"expires" validate:"required"`
}

type AccessUpdate struct {
	ProductID  string    `firestore:"dataproduct_id" json:"dataproduct_id"`
	Author     string    `firestore:"author" json:"author"`
	Subject    string    `firestore:"subject" json:"subject"`
	Action     string    `firestore:"action" json:"action"`
	UpdateTime time.Time `firestore:"time" json:"time"`
	Expires    time.Time `firestore:"expires" json:"expires"`
}

func (a *api) getAccessUpdatesForProduct(w http.ResponseWriter, r *http.Request) {
	updates := a.client.Collection(a.config.Firestore.AccessUpdatesCollection)
	productID := chi.URLParam(r, "productID")

	updateResponse := make([]AccessUpdate, 0)

	query := updates.Where("dataproduct_id", "==", productID).OrderBy("time", firestore.Desc)
	iter := query.Documents(r.Context())
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

		var update AccessUpdate
		if err := document.DataTo(&update); err != nil {
			log.Errorf("Deserializing firestore document: %v", err)
			respondf(w, http.StatusInternalServerError, "unable to deserialize update\n")
			return
		}

		updateResponse = append(updateResponse, update)
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(updateResponse); err != nil {
		log.Errorf("Serializing updateResponses: %v", err)
		respondf(w, http.StatusInternalServerError, "unable to serialize updateResponses\n")
		return
	}
}

func (a *api) getAccessForProduct(w http.ResponseWriter, r *http.Request) {
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

	dpr, err := DocumentToProductResponse(document)
	if err != nil {
		log.Errorf("Deserializing firestore document: %v", err)
		respondf(w, http.StatusInternalServerError, "unable to deserialize document\n")
		return
	}

	var response AccessResponse
	response.Access = dpr.DataProduct.Access
	response.ID = dpr.ID
	response.Created = dpr.Created
	response.Updated = dpr.Updated

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Errorf("Serializing access response: %v", err)
		respondf(w, http.StatusInternalServerError, "unable to serialize access response\n")
		return
	}
}

func (a *api) removeAccessForProduct(w http.ResponseWriter, r *http.Request) {
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

	dpr, err := DocumentToProductResponse(document)
	if err != nil {
		log.Errorf("Deserializing firestore document: %v", err)
		respondf(w, http.StatusInternalServerError, "unable to deserialize document\n")
		return
	}

	var accessSubject AccessSubject
	if err := json.NewDecoder(r.Body).Decode(&accessSubject); err != nil {
		log.Errorf("Deserializing request document: %v", err)
		respondf(w, http.StatusBadRequest, "unable to deserialize request document\n")
		return
	}

	requester := r.Context().Value("preferred_username").(string)

	if dpr.DataProduct.Owner == requester || accessSubject.Subject == requester {
		_, ok := dpr.DataProduct.Access[accessSubject.Subject]
		if !ok {
			log.Errorf("Requested subject does have an access entry")
			respondf(w, http.StatusBadRequest, "requested subject does not have an access entry")
			return
		}

		delete(dpr.DataProduct.Access, accessSubject.Subject)
		documentRef.Update(r.Context(), []firestore.Update{{
			Path:  "access",
			Value: dpr.DataProduct.Access,
		}})
		iam.RemoveDatastoreAccess(r.Context(), dpr.DataProduct.Datastore[0], accessSubject.Subject)

		update := AccessUpdate{
			Subject:    accessSubject.Subject,
			Action:     Delete,
			ProductID:  dpr.ID,
			UpdateTime: time.Now(),
			Author:     requester,
		}

		UpdateHistory(r.Context(), a.client, a.config.Firestore.AccessUpdatesCollection, update)
		w.WriteHeader(http.StatusNoContent)
		return
	}

	log.Errorf("Requester is not authorized to make changes to this rule: product id: %v, requester: %v, subject: %v", dpr.ID, requester, accessSubject.Subject)
	respondf(w, http.StatusUnauthorized, "you are unauthorized to make changes to this access rule")

}

func (a *api) grantAccessForProduct(w http.ResponseWriter, r *http.Request) {
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

	dpr, err := DocumentToProductResponse(document)
	if err != nil {
		log.Errorf("Deserializing firestore document: %v", err)
		respondf(w, http.StatusInternalServerError, "unable to deserialize document\n")
		return
	}

	var accessSubject AccessSubject
	if err := json.NewDecoder(r.Body).Decode(&accessSubject); err != nil {
		log.Errorf("Deserializing request document: %v", err)
		respondf(w, http.StatusBadRequest, "unable to deserialize request document\n")
		return
	}

	if err := a.validate.Struct(accessSubject); err != nil {
		log.Errorf("Validating request document: %v", err)
		respondf(w, http.StatusBadRequest, "unable to validate request document\n")
		return
	}

	requester := r.Context().Value("preferred_username").(string)

	if dpr.DataProduct.Owner != requester {
		log.Errorf("Requester is not authorized to make changes to this rule: product id: %v, requester: %v, subject: %v", dpr.ID, requester, accessSubject.Subject)
		respondf(w, http.StatusUnauthorized, "you are unauthorized to make changes to this access rule")
		return
	}

	dpr.DataProduct.Access[accessSubject.Subject] = accessSubject.Expires
	documentRef.Update(r.Context(), []firestore.Update{{
		Path:  "access",
		Value: dpr.DataProduct.Access,
	}})
	iam.UpdateDatastoreAccess(r.Context(), dpr.DataProduct.Datastore[0], dpr.DataProduct.Access)

	update := AccessUpdate{
		Subject:    accessSubject.Subject,
		Action:     Grant,
		ProductID:  dpr.ID,
		Expires:    accessSubject.Expires,
		UpdateTime: time.Now(),
		Author:     requester,
	}
	UpdateHistory(r.Context(), a.client, a.config.Firestore.AccessUpdatesCollection, update)
	w.WriteHeader(http.StatusNoContent)
}

func UpdateHistory(ctx context.Context, client *firestore.Client, collectionName string, update AccessUpdate) {
	updates := client.Collection(collectionName)
	updates.Add(ctx, update)
}
