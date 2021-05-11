package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/nais/dp/backend/auth"
	"google.golang.org/api/iterator"

	"cloud.google.com/go/firestore"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/nais/dp/backend/config"
	"github.com/nais/dp/backend/middleware"
	"golang.org/x/oauth2"
	"gopkg.in/go-playground/validator.v9"

	log "github.com/sirupsen/logrus"
)

const (
	BucketType   = "bucket"
	BigQueryType = "bigquery"
)

type api struct {
	client   *firestore.Client
	validate *validator.Validate
	config   config.Config
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

func New(client *firestore.Client, config config.Config) chi.Router {
	api := api{client, validator.New(), config}

	latencyHistBuckets := []float64{.001, .005, .01, .025, .05, .1, .5, 1, 3, 5}
	prometheusMiddleware := middleware.PrometheusMiddleware("backend", latencyHistBuckets...)
	prometheusMiddleware.Initialize("/api/v1/", http.MethodGet, http.StatusOK)

	r := chi.NewRouter()

	r.Use(prometheusMiddleware.Handler())
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	}))

	r.Route("/api/v1", func(r chi.Router) {
		// requires valid access token
		r.Group(func(r chi.Router) {
			r.Use(middleware.JWTValidatorMiddleware(auth.KeyDiscoveryURL(config.OAuth2TenantID), config.OAuth2ClientID, config.DevMode))
			r.Post("/dataproducts", api.createDataproduct)
			r.Put("/dataproducts/{productID}", api.updateDataproduct)
			r.Delete("/dataproducts/{productID}", api.deleteDataproduct)
			r.Get("/teams", api.getTeamsForUser)
		})

		r.Get("/dataproducts", api.dataproducts)
		r.Get("/dataproducts/{productID}", api.getDataproduct)

	})

	r.Get("/callback", api.callback)

	return r
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

		var dpr DataProductResponse
		var dp DataProduct

		err = document.DataTo(&dp)

		if dp.Access == nil {
			dp.Access = make([]*AccessEntry, 0)
		}

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
	dpc := a.client.Collection(a.config.FirestoreCollection)
	var dp DataProduct

	if err := json.NewDecoder(r.Body).Decode(&dp); err != nil {
		log.Errorf("Deserializing document: %v", err)
		respondf(w, http.StatusBadRequest, "unable to deserialize document\n")
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
		log.Errorf("Getting document: %v", err)
		respondf(w, http.StatusNotFound, "unable to get document\n")
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
		log.Errorf("Deserializing document: %v", err)
		respondf(w, http.StatusBadRequest, "unable to deserialize document\n")
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
		log.Errorf("Updating document: %v", err)
		respondf(w, http.StatusBadRequest, "unable to update document\n")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (a *api) getDataproduct(w http.ResponseWriter, r *http.Request) {
	dpc := a.client.Collection(a.config.FirestoreCollection)
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
	dpc := a.client.Collection(a.config.FirestoreCollection)
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

func (a *api) createUpdates(dp DataProduct, existingDp DataProduct) ([]firestore.Update, error) {
	var updates []firestore.Update
	newAccess := existingDp.Access

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
	if len(dp.Owner) > 0 {
		updates = append(updates, firestore.Update{
			Path:  "owner",
			Value: dp.Owner,
		})
	}
	if len(dp.Access) > 0 {
		for _, access := range dp.Access {
			if errs := a.validate.Struct(access); errs != nil {
				return nil, errs
			}
		}
		newAccess = append(newAccess, dp.Access...)
		updates = append(updates, firestore.Update{
			Path:  "access",
			Value: newAccess,
		})
	}

	return updates, nil
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
	teamsMap, err := fetchAllTeams(a.config.TeamsURL)
	if err != nil {
		log.Errorf("fetching teams: %v", err)
		respondf(w, http.StatusInternalServerError, "unable to fetch teams\n")
		return
	}
	var teams []string
	for _, uuid := range r.Context().Value("groups").([]string) {
		if _, found := teamsMap[uuid]; found {
			teams = append(teams, teamsMap[uuid])
		}
	}

	if err := json.NewEncoder(w).Encode(teams); err != nil {
		log.Errorf("encoding teams response: %v", err)
		respondf(w, http.StatusInternalServerError, "unable to get teams for user\n")
		return
	}
}

func fetchAllTeams(url string) (map[string]string, error) {
	var teams map[string]string
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if err := json.NewDecoder(response.Body).Decode(&teams); err != nil {
		return nil, err
	}
	return teams, nil
}

func ValidateDatastore(store map[string]string) error {
	datastoreType := store["type"]
	if len(datastoreType) == 0 {
		return fmt.Errorf("no type defined")
	}

	switch datastoreType {
	case BucketType:
		return hasKeys(store, "project_id", "bucket_id")
	case BigQueryType:
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
