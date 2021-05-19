package api

import (
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/nais/dp/backend/auth"
	"github.com/nais/dp/backend/config"
	"github.com/nais/dp/backend/middleware"

	"gopkg.in/go-playground/validator.v9"
)

func New(client *firestore.Client, config config.Config, teamUUIDs map[string]string) chi.Router {
	api := api{
		client:    client,
		validate:  validator.New(),
		config:    config,
		teamUUIDs: teamUUIDs,
	}

	azureGroups := auth.AzureGroups{
		Cache:  make(map[string]auth.CacheEntry),
		Client: http.DefaultClient,
		Config: config,
	}

	latencyHistBuckets := []float64{.001, .005, .01, .025, .05, .1, .5, 1, 3, 5}
	prometheusMiddleware := middleware.PrometheusMiddleware("backend", latencyHistBuckets...)
	prometheusMiddleware.Initialize("/api/v1/", http.MethodGet, http.StatusOK)
	authenticatorMiddleware := middleware.JWTValidatorMiddleware(auth.KeyDiscoveryURL(config.OAuth2.TenantID), config.OAuth2.ClientID, config.DevMode, azureGroups)

	r := chi.NewRouter()

	r.Use(prometheusMiddleware.Handler())
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowCredentials: true,
	}))

	r.Get("/oauth2/callback", api.callback)
	r.Get("/login", api.login)

	r.Route("/api/v1", func(r chi.Router) {
		// requires valid access token
		r.Group(func(r chi.Router) {
			r.Use(authenticatorMiddleware)
			r.Get("/userinfo", api.userInfo)
		})
		r.Route("/dataproducts", func(r chi.Router) {
			r.Group(func(r chi.Router) {
				r.Use(authenticatorMiddleware)
				r.Post("/", api.createDataproduct)
				r.Put("/{productID}", api.updateDataproduct)
				r.Delete("/{productID}", api.deleteDataproduct)
			})
			r.Get("/", api.dataproducts)
			r.Get("/{productID}", api.getDataproduct)
		})
		r.Route("/access", func(r chi.Router) {
			r.Group(func(r chi.Router) {
				r.Use(authenticatorMiddleware)
				r.Delete("/{productID}", api.removeAccessForProduct)
				r.Post("/{productID}", api.grantAccessForProduct)
				r.Get("/{productID}/history", api.getAccessUpdatesForProduct)
			})
			r.Get("/{productID}", api.getAccessForProduct)
		})
	})

	return r
}
