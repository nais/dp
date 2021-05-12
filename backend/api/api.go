package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/nais/dp/backend/auth"
	"google.golang.org/api/iterator"

	"github.com/go-chi/chi"
	"golang.org/x/oauth2"

	log "github.com/sirupsen/logrus"
)

const (
	BucketType   = "bucket"
	BigQueryType = "bigquery"
)

func (a *api) getDataproduct(w http.ResponseWriter, r *http.Request) {
	dpc := a.client.Collection(a.config.FirestoreCollection)
	articleID := chi.URLParam(r, "productID")
	documentRef := dpc.Doc(articleID)

	document, err := documentRef.Get(r.Context())
	if err != nil {
		log.Errorf("Getting firestore document: %v", err)
		respondf(w, http.StatusBadRequest, "unable to get document\n")
		return
	}

	dpr, err := documentToProduct(document)
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
	dpc := a.client.Collection(a.config.FirestoreCollection)
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

		dpr, err := documentToProduct(document)
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
	dpc := a.client.Collection(a.config.FirestoreCollection)
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
	dpc := a.client.Collection(a.config.FirestoreCollection)
	articleID := chi.URLParam(r, "productID")
	documentRef := dpc.Doc(articleID)
	document, err := documentRef.Get(r.Context())
	if err != nil {
		log.Errorf("Getting firestore document: %v", err)
		respondf(w, http.StatusNotFound, "unable to get firestore document\n")
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

	updates, err := a.createUpdates(dp, firebaseDp)
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

func (a *api) deleteDataproduct(w http.ResponseWriter, r *http.Request) {
	dpc := a.client.Collection(a.config.FirestoreCollection)
	articleID := chi.URLParam(r, "productID")
	documentRef := dpc.Doc(articleID)

	if _, err := documentRef.Delete(r.Context()); err != nil {
		log.Errorf("Deleting firestore document: %v", err)
		respondf(w, http.StatusBadRequest, "unable to delete firestore document\n")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (a *api) callback(w http.ResponseWriter, r *http.Request) {
	cfg := auth.CreateOAuth2Config(a.config)

	state := "veryrandomstring"
	consentUrl := cfg.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.SetAuthURLParam("redirect_uri", "http://localhost:8080/callback"))
	fmt.Println(consentUrl)

	code := r.URL.Query().Get("code")
	if len(code) == 0 {
		respondf(w, http.StatusForbidden, "No code in query params")
		return
	}

	tokens, err := cfg.Exchange(r.Context(), code)
	if err != nil {
		log.Errorf("Exchanging authorization code for tokens: %v", err)
		respondf(w, http.StatusForbidden, "uh oh")
		return
	}

	//w.Header().Set("Set-Cookie", fmt.Sprintf("access_token=%v;HttpOnly;Secure;Max-Age=86400;Domain=%v", tokens.AccessToken, "dp.dev.intern.nav.no"))
	w.Header().Set("Set-Cookie", fmt.Sprintf("jwt=%v;HttpOnly;Secure;Max-Age=86400", tokens.AccessToken))
	w.WriteHeader(http.StatusOK)
}

func (a *api) getTeamsForUser(w http.ResponseWriter, r *http.Request) {
	var teams []string
	for _, uuid := range r.Context().Value("groups").([]string) {
		if _, found := a.teamUUIDs[uuid]; found {
			teams = append(teams, a.teamUUIDs[uuid])
		}
	}

	log.Infof("User groups: %v, group names: %v", r.Context().Value("groups"), teams)

	if err := json.NewEncoder(w).Encode(teams); err != nil {
		log.Errorf("Serializing teams response: %v", err)
		respondf(w, http.StatusInternalServerError, "unable to serialize teams for user\n")
		return
	}
}
