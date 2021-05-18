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

	r := chi.NewRouter()

	r.Use(prometheusMiddleware.Handler())
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowCredentials: true,
	}))

	r.Route("/api/v1", func(r chi.Router) {
		// requires valid access token
		r.Group(func(r chi.Router) {
			r.Use(middleware.JWTValidatorMiddleware(auth.KeyDiscoveryURL(config.OAuth2TenantID), config.OAuth2ClientID, config.DevMode, azureGroups))
			r.Post("/dataproducts", api.createDataproduct)
			r.Put("/dataproducts/{productID}", api.updateDataproduct)
			r.Delete("/dataproducts/{productID}", api.deleteDataproduct)
			r.Get("/userinfo", api.userInfo)
		})

		r.Get("/dataproducts", api.dataproducts)
		r.Get("/dataproducts/{productID}", api.getDataproduct)
	})

	r.Get("/oauth2/callback", api.callback)
	r.Get("/login", api.login)

	return r
}
