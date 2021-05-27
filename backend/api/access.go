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
	firestore2 "github.com/nais/dp/backend/firestore"
	"github.com/nais/dp/backend/iam"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/iterator"
)

const (
	UserType           = "user"
	ServiceAccountType = "serviceAccount"
)

type AccessSubject struct {
	Subject string    `json:"subject" validate:"required"`
	Type    string    `json:"type" validate:"required"`
	Expires time.Time `json:"expires"`
}

func (a *api) getAccessUpdatesForProduct(w http.ResponseWriter, r *http.Request) {
	updates := a.client.Collection(a.config.Firestore.AccessUpdatesCollection)
	productID := chi.URLParam(r, "productID")

	updateResponse := make([]firestore2.AccessUpdate, 0)

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

		var update firestore2.AccessUpdate
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

func (a *api) removeProductAccess(w http.ResponseWriter, r *http.Request) {
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

	var accessSubject AccessSubject
	if err := json.NewDecoder(r.Body).Decode(&accessSubject); err != nil {
		log.Errorf("Deserializing request document: %v", err)
		respondf(w, http.StatusBadRequest, "unable to deserialize request document\n")
		return
	}

	var subject string
	switch accessSubject.Type {
	case UserType:
		subject = "user:" + accessSubject.Subject
	case ServiceAccountType:
		subject = "serviceAccount:" + accessSubject.Subject
	default:
		{
			log.Errorf("Invalid AccessSubject.Type: %v", accessSubject.Type)
			respondf(w, http.StatusBadRequest, "invalid AccessSubject.Type\n")
			return
		}
	}

	requester := r.Context().Value("preferred_username").(string)
	requesterMember := r.Context().Value("member_name").(string)
	requesterGroups := r.Context().Value("teams").([]string)

	if contains(requesterGroups, dpr.Dataproduct.Team) || accessSubject.Subject == requesterMember {
		_, ok := dpr.Dataproduct.Access[subject]
		if !ok {
			log.Errorf("Requested subject does have an access entry")
			respondf(w, http.StatusBadRequest, "requested subject does not have an access entry")
			return
		}

		delete(dpr.Dataproduct.Access, subject)
		documentRef.Update(r.Context(), []firestore.Update{{
			Path:  "access",
			Value: dpr.Dataproduct.Access,
		}})
		iam.RemoveDatastoreAccess(r.Context(), dpr.Dataproduct.Datastore[0], subject)

		update := firestore2.Delete(requester, dpr.ID, accessSubject.Subject)
		if err := a.firestore.AddAccessUpdate(r.Context(), update); err != nil {
			log.Errorf("Adding access update: %v", err)
		}

		w.WriteHeader(http.StatusNoContent)
		return
	}

	log.Errorf("Requester is not authorized to make changes to this rule: product id: %v, requester: %v, subject: %v", dpr.ID, requester, accessSubject.Subject)
	respondf(w, http.StatusUnauthorized, "you are unauthorized to make changes to this access rule")
}

func (a *api) grantProductAccess(w http.ResponseWriter, r *http.Request) {
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

	if accessSubject.Expires.Before(time.Now()) && !accessSubject.Expires.IsZero() {
		log.Errorf("Invalid AccessSubject.Expires: %v is already an expired time", accessSubject.Expires)
		respondf(w, http.StatusBadRequest, "invalid AccessSubject.Expires\n")
		return
	}

	var subject string
	switch accessSubject.Type {
	case UserType:
		subject = "user:" + accessSubject.Subject
	case ServiceAccountType:
		subject = "serviceAccount:" + accessSubject.Subject
	default:
		{
			log.Errorf("Invalid AccessSubject.Type: %v", accessSubject.Type)
			respondf(w, http.StatusBadRequest, "invalid AccessSubject.Type\n")
			return
		}
	}

	requester := r.Context().Value("preferred_username").(string)

	dpr.Dataproduct.Access[subject] = accessSubject.Expires
	documentRef.Update(r.Context(), []firestore.Update{{
		Path:  "access",
		Value: dpr.Dataproduct.Access,
	}})
	iam.UpdateDatastoreAccess(r.Context(), dpr.Dataproduct.Datastore[0], dpr.Dataproduct.Access)

	update := firestore2.Grant(requester, dpr.ID, accessSubject.Subject, accessSubject.Expires)
	UpdateHistory(r.Context(), a.client, a.config.Firestore.AccessUpdatesCollection, update)
	w.WriteHeader(http.StatusNoContent)
}

func UpdateHistory(ctx context.Context, client *firestore.Client, collectionName string, update firestore2.AccessUpdate) {
	updates := client.Collection(collectionName)
	updates.Add(ctx, update)
}

func documentToProductResponse(d *firestore.DocumentSnapshot) (firestore2.DataproductResponse, error) {
	var dpr firestore2.DataproductResponse
	var dp firestore2.Dataproduct

	if err := d.DataTo(&dp); err != nil {
		return dpr, err
	}

	if dp.Access == nil {
		dp.Access = make(map[string]time.Time)
	}
	dpr.ID = d.Ref.ID
	dpr.Updated = d.UpdateTime
	dpr.Created = d.CreateTime
	dpr.Dataproduct = dp

	return dpr, nil
}
